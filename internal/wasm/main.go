//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"

	keyrecovery "github.com/0x3327/curvy-social-recovery/key_recovery"
	BN254_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	SECP256K1_fr "github.com/consensys/gnark-crypto/ecc/secp256k1/fr"
)

func jsSplit(this js.Value, args []js.Value) (any, error) {
	t := args[0].Int()
	n := args[1].Int()
	skStr := args[2].String()
	vkStr := args[3].String()

	shares, err := keyrecovery.Split(t, n, skStr, vkStr)
	if err != nil {
		return nil, err
	}

	ret := make([]any, 0, len(shares))

	for _, share := range shares {
		shareMap := make(map[string]any)
		shareMap["x"] = share.Point
		shareMap["skEval"] = share.SpendingEval
		shareMap["vkEval"] = share.ViewingEval

		ret = append(ret, shareMap)
	}

	return ret, nil
}

func jsRec(this js.Value, args []js.Value) (any, error) {
	t := args[0].Int()
	shareMapSlice := args[1]
	nOfShares := shareMapSlice.Length()

	shares := make([]keyrecovery.Share, 0, nOfShares)

	for i := 0; i < nOfShares; i++ {
		point := shareMapSlice.Index(i).Get("x").String()
		skEval := shareMapSlice.Index(i).Get("skEval").String()
		vkEval := shareMapSlice.Index(i).Get("vkEval").String()
		shares = append(shares, keyrecovery.Share{
			Point:        point,
			SpendingEval: skEval,
			ViewingEval:  vkEval,
		})
	}

	skStr, vkStr, err := keyrecovery.Recover(t, shares)

	if err != nil {
		return nil, err
	}

	ret := make(map[string]any)
	ret["k"] = skStr
	ret["v"] = vkStr

	return ret, nil
}

func jsGenRandomKeys(this js.Value, args []js.Value) (any, error) {
	var spendingKey SECP256K1_fr.Element
	var viewingKey BN254_fr.Element

	_, err := spendingKey.SetRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate random spending key: %v\n", err)
	}
	_, err = viewingKey.SetRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate random viewing key: %v\n", err)
	}

	spendingKeyStr := spendingKey.Text(16)
	viewingKeyStr := viewingKey.Text(16)

	ret := make(map[string]any)
	ret["k"] = spendingKeyStr
	ret["v"] = viewingKeyStr

	return ret, nil
}

func promiseWrapper(fn func(js.Value, []js.Value) (any, error)) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) any {

		handler := js.FuncOf(func(_ js.Value, handlerArgs []js.Value) any {
			resolve := handlerArgs[0]
			reject := handlerArgs[1]

			go func() {
				defer func() {
					if err := recover(); err != nil {
						errConstructor := js.Global().Get("Error")
						errObject := errConstructor.New(fmt.Sprintf("%v", err))
						reject.Invoke(errObject)
					}
				}()
				res, err := fn(this, args)
				if err != nil {
					errConstructor := js.Global().Get("Error")
					errObject := errConstructor.New(err.Error())
					reject.Invoke(errObject)
				} else {
					resolve.Invoke(js.ValueOf(res))
				}
			}()

			return nil
		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	}
}

func main() {
	wait := make(chan struct{}, 0)
	js.Global().Set("goSplit", js.FuncOf(promiseWrapper(jsSplit)))
	js.Global().Set("goRecover", js.FuncOf(promiseWrapper(jsRec)))
	js.Global().Set("goGenRandomKeys", js.FuncOf(promiseWrapper(jsGenRandomKeys)))
	<-wait
}

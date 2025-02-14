package keyrecovery

import (
	"errors"

	BN254_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	SECP256K1_fr "github.com/consensys/gnark-crypto/ecc/secp256k1/fr"
)

func Split(threshold, nOfShares int, spendingKey, viewingKey string) ([]Share, error) {

	if threshold > nOfShares {
		return nil, errors.New("threshold t should be less than or equal to number of shares n")
	}

	// Generate the polynomials for the spending key and viewing key
	skPolynomial := make([]SECP256K1_fr.Element, threshold)
	vkPolynomial := make([]BN254_fr.Element, threshold)

	// Generate random coefficients for the polynomials
	for i := 1; i < threshold; i++ {
		_, errs := skPolynomial[i].SetRandom()
		_, errv := vkPolynomial[i].SetRandom()

		if errs != nil {
			return nil, errs
		}
		if errv != nil {
			return nil, errv
		}
	}

	// Set the free coefficients of the polynomials to the spending key and viewing key
	_, errs := skPolynomial[0].SetString("0x" + spendingKey)
	_, errv := vkPolynomial[0].SetString("0x" + viewingKey)

	if errs != nil {
		return nil, errs
	}
	if errv != nil {
		return nil, errv
	}

	// Generate the shares
	shares := make([]Share, nOfShares)
	for i := range shares {
		// Define x coordinate for the share as field element
		var xsk SECP256K1_fr.Element
		var xvk BN254_fr.Element
		xsk.SetUint64(uint64(i + 1))
		xvk.SetUint64(uint64(i + 1))

		// Evaluate the polynomials at the x coordinate
		spendingEval := evalSpending(skPolynomial, &xsk)
		viewingEval := evalViewing(vkPolynomial, &xvk)

		// Set the share
		shares[i].Point = xsk.Text(16)
		shares[i].SpendingEval = spendingEval.Text(16)
		shares[i].ViewingEval = viewingEval.Text(16)
	}

	return shares, nil
}

// https://github.com/Consensys/gnark-crypto/blob/7c56ee003026d11987e8f965bf2674b4ca052ea8/ecc/bn254/fr/polynomial/polynomial.go#L25

func evalSpending(p []SECP256K1_fr.Element, x *SECP256K1_fr.Element) SECP256K1_fr.Element {
	res := p[len(p)-1]
	for i := len(p) - 2; i >= 0; i-- {
		res.Mul(&res, x)
		res.Add(&res, &(p)[i])
	}
	return res
}

func evalViewing(p []BN254_fr.Element, x *BN254_fr.Element) BN254_fr.Element {
	res := p[len(p)-1]
	for i := len(p) - 2; i >= 0; i-- {
		res.Mul(&res, x)
		res.Add(&res, &(p)[i])
	}
	return res
}

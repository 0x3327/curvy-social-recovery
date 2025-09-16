package keyrecovery

import (
	"errors"

	BN254_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	SECP256K1_fr "github.com/consensys/gnark-crypto/ecc/secp256k1/fr"
)

// Recover returns the recovered spending key, viewing key and error encountered, if any.
//
// Keys returned are hex strings without the leading 0x
func Recover(threshold int, shares []Share) (string, string, error) {

	// Check if the number of shares is less than the threshold
	if len(shares) < threshold {
		return "", "", &NumberOfSharesLessThanThreshold{
			N: len(shares),
			T: threshold,
		}
	}

	if valid, idx := pointsUnique(shares); !valid {
		return "", "", &DuplicatePointInSharesError{
			Idx:   idx,
			Point: shares[idx].Point,
		}
	}

	points, err := extractAllPointsFromShares(shares)
	if err != nil {
		return "", "", err
	}

	// We recover with the first t=threshold points
	skStr, vkStr, err := recoverFromPoints(points[0:threshold])
	if err != nil {
		return "", "", err
	}

	if len(points) == threshold {
		return skStr, vkStr, nil
	}

	// We recover with the last t=threshold points to verify
	skStr2, vkStr2, err := recoverFromPoints(points[len(points)-threshold:])
	if err != nil {
		return "", "", err
	}

	if skStr != skStr2 || vkStr != vkStr2 {
		return "", "", &RecoveredKeysDoNotMatchError{}
	}

	return skStr, vkStr, nil
}

func recoverFromPoints(points []point) (string, string, error) {

	// Define the spending key and viewing key as field elements for calculations
	var sk SECP256K1_fr.Element
	var vk BN254_fr.Element

	// Initialize the spending key and viewing key as zero
	sk.SetZero()
	vk.SetZero()

	// Recover the spending key and viewing key from the shares
	//  		n             n
	// Result = Σ  y_i   *    Π ( x_j / (x_j - x_i) )
	// 		   i=1           j=1
	// 						j ≠ i

	for i := range points {
		// Current term = y_i
		var lsk SECP256K1_fr.Element
		var lvk BN254_fr.Element
		lsk.Set(&points[i].yS)
		lvk.Set(&points[i].yV)

		// Denominators for the Lagrange basis Li
		denominatorS := SECP256K1_fr.One()
		denominatorV := BN254_fr.One()

		// Apply Li to the current term
		for j := range points {
			if i == j {
				continue
			}

			// current term * xj
			lsk.Mul(&lsk, &points[j].xS)
			lvk.Mul(&lvk, &points[j].xV)

			// xj - xi
			var diffS SECP256K1_fr.Element
			var diffV BN254_fr.Element
			diffS.Sub(&points[j].xS, &points[i].xS)
			diffV.Sub(&points[j].xV, &points[i].xV)

			// denominator *= (xj - xi)
			denominatorS.Mul(&denominatorS, &diffS)
			denominatorV.Mul(&denominatorV, &diffV)
		}

		// Divide by the denominator
		lsk.Div(&lsk, &denominatorS)
		lvk.Div(&lvk, &denominatorV)

		// Add the current term to the spending key and viewing key
		sk.Add(&sk, &lsk)
		vk.Add(&vk, &lvk)
	}

	return sk.Text(16), vk.Text(16), nil
}

func pointsUnique(shares []Share) (bool, int) {
	for i := range shares {
		for j := i + 1; j < len(shares); j++ {
			if shares[i].Point == shares[j].Point {
				return false, i
			}
		}
	}
	return true, 0
}

func extractAllPointsFromShares(shares []Share) ([]point, error) {
	return extractPointsFromShares(shares, len(shares))
}

// extractPointFromShares returns no more than maxNOfPoints distinct points from shares
func extractPointsFromShares(shares []Share, maxNOfPoints int) ([]point, error) {
	if len(shares) == 0 {
		return nil, errors.New("no shares provided")
	}

	points := make([]point, 0, maxNOfPoints)

	for i, share := range shares {
		if i >= maxNOfPoints {
			break
		}

		var p point
		err := p.fromShare(share)
		if err != nil {
			return nil, err
		}

		points = append(points, p)
	}

	return points, nil
}

type point struct {
	xS SECP256K1_fr.Element
	xV BN254_fr.Element
	yS SECP256K1_fr.Element
	yV BN254_fr.Element
}

func (p *point) fromShare(s Share) error {
	_, errX1 := p.xS.SetString("0x" + s.Point)
	_, errX2 := p.xV.SetString("0x" + s.Point)
	_, errY1 := p.yS.SetString("0x" + s.SpendingEval)
	_, errY2 := p.yV.SetString("0x" + s.ViewingEval)

	if errX1 != nil {
		return errX1
	}
	if errX2 != nil {
		return errX2
	}
	if errY1 != nil {
		return errY1
	}
	if errY2 != nil {
		return errY2
	}

	return nil
}

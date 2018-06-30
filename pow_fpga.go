// +build cgo
// +build linux

package giota

// #cgo CFLAGS: -Wall -fsigned-char

/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <error.h>
#include <stdint.h>
#include <string.h>
#include <sys/mman.h>
#include <fcntl.h>

#define NONCE_LEN                   81
#define MWM_MASK_REG_OFFSET         3

long long int pworkFPGA(signed char itrits[], int itrits_len, int mwm, signed char nonce[])
{
    FILE *ctrl_fd = 0;
    FILE *in_fd = 0;
    FILE *out_fd = 0;

    int result;

    ctrl_fd = fopen("/dev/cpow-ctrl", "r+");

    if(ctrl_fd == NULL) {
        return -1;
    }

    in_fd = fopen("/dev/cpow-idata", "wb");

    if(in_fd == NULL) {
        fclose(ctrl_fd);
	    return -1;
    }

    out_fd = fopen("/dev/cpow-odata", "rb");

    if(out_fd == NULL) {
 	    fclose(ctrl_fd);
	    fclose(in_fd);
        return -1;
    }

    fwrite(itrits, 1, itrits_len, in_fd);
    fflush(in_fd);

    fwrite(&mwm, 1, 1, ctrl_fd);
    fread(&result, sizeof(result), 1, ctrl_fd);
    fflush(ctrl_fd);

    fread(nonce, 1, NONCE_LEN, out_fd);

    fclose(in_fd);
    fclose(out_fd);
    fclose(ctrl_fd);

    return 0;
}

*/
import "C"
import (
	"errors"
	"unsafe"
)

func PowFPGA(trytes Trytes, mwm int) (Trytes, error) {

	var (
		isPOW bool
		ilen int
		result Trytes
	)

	if isPOW == true {
		return "", errors.New("pow is already running, stopped")
	}

	if trytes == "" {
		return "", errors.New("invalid trytes")
	}

	isPOW = true

	tr := trytes.Trits()
	ilen = len(tr)

	nonce := make(Trits, NonceTrinarySize)

	C.pworkFPGA((*C.schar)(unsafe.Pointer(&tr[0])), C.int(ilen) , C.int(mwm), (*C.schar)(unsafe.Pointer(&nonce[0])) )

	result = nonce.Trytes()

	isPOW = false

	return result, nil
}


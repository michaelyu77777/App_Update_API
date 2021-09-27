package tlss

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/sirupsen/logrus"
	"leapsy.com/packages/logings"
)

// CreateCertAndKeyPEMFiles - 產生證書和鑰匙PEM檔案
/**
 * @param  string certPEMFileName  證書PEM檔名
 * @param  string privateKeyPEMFileName 私鑰PEM檔名
 */
func CreateCertAndKeyPEMFiles(certPEMFileName, privateKeyPEMFileName string) {

	privateKeyPointer, ecdsaGenerateKeyError := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	logings.SendLog(
		[]string{`產生私鑰 `},
		[]interface{}{},
		ecdsaGenerateKeyError,
		logrus.PanicLevel,
	)

	serialNumber, randIntError := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))

	logings.SendLog(
		[]string{`產生序號 `},
		[]interface{}{},
		randIntError,
		logrus.PanicLevel,
	)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{`Leapsy`},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now(),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, x509CreateCertificateError := x509.CreateCertificate(rand.Reader, &template, &template, &privateKeyPointer.PublicKey, privateKeyPointer)

	logings.SendLog(
		[]string{`產生證書 `},
		[]interface{}{},
		x509CreateCertificateError,
		logrus.PanicLevel,
	)

	certPEMBytes := pem.EncodeToMemory(&pem.Block{Type: `CERTIFICATE`, Bytes: certBytes})

	var pemEncodeToMemoryError error

	if certPEMBytes == nil {
		pemEncodeToMemoryError = fmt.Errorf(`加密成PEM結果為空`)
	}

	logings.SendLog(
		[]string{`將 證書 '%s' 加密為 PEM '%s' `},
		[]interface{}{string(certBytes), string(certPEMBytes)},
		pemEncodeToMemoryError,
		logrus.PanicLevel,
	)

	ioutilWriteCertPEMFileError := ioutil.WriteFile(certPEMFileName, certPEMBytes, 0644)

	logings.SendLog(
		[]string{`將 證書PEM '%s' 寫入檔案 '%s' `},
		[]interface{}{string(certPEMBytes), certPEMFileName},
		ioutilWriteCertPEMFileError,
		logrus.PanicLevel,
	)

	privateKeyBytes, x509MarshalPKCS8PrivateKeyError := x509.MarshalPKCS8PrivateKey(privateKeyPointer)

	logings.SendLog(
		[]string{`將 私鑰物件 %+v 整理成私鑰 '%s' `},
		[]interface{}{*privateKeyPointer, string(privateKeyBytes)},
		x509MarshalPKCS8PrivateKeyError,
		logrus.PanicLevel,
	)

	privateKeyPEMBytes := pem.EncodeToMemory(&pem.Block{Type: `PRIVATE KEY`, Bytes: privateKeyBytes})

	if privateKeyPEMBytes == nil {
		pemEncodeToMemoryError = fmt.Errorf(`加密成PEM結果為空`)
	}

	logings.SendLog(
		[]string{`將 私鑰 '%s' 加密為 PEM '%s' `},
		[]interface{}{string(privateKeyBytes), string(privateKeyPEMBytes)},
		pemEncodeToMemoryError,
		logrus.PanicLevel,
	)

	ioutilWritePrivateKeyPEMFileError := ioutil.WriteFile(privateKeyPEMFileName, privateKeyPEMBytes, 0600)

	logings.SendLog(
		[]string{`將 私鑰PEM '%s' 寫入檔案 '%s' `},
		[]interface{}{string(privateKeyPEMBytes), privateKeyPEMFileName},
		ioutilWritePrivateKeyPEMFileError,
		logrus.PanicLevel,
	)

}

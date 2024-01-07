/*
 * This example showcases how to create appearance fields for digital signature.
 * The file is signed using a private/public key pair.
 *
 * $ ./pdf_image_sign_appearance <INPUT_PDF_PATH> <IMAGE_FILE> <WATERMARK_IMAGE_FILE> <OUTPUT_PDF_PATH>
 */
 package unipdflib

 import (
	 "crypto/rand"
	 "crypto/rsa"
	 "crypto/x509"
	 "crypto/x509/pkix"
	 "fmt"
	 "image"
	 "log"
	 "math/big"
	 "os"
	 "time"
	 "github.com/unidoc/unipdf/v3/annotator"
	 "github.com/unidoc/unipdf/v3/common/license"
	 "github.com/unidoc/unipdf/v3/core"
	 "github.com/unidoc/unipdf/v3/model"
	 "github.com/unidoc/unipdf/v3/model/sighandler"
 )
 
 var now = time.Now()
 

 func SignWithImage(inputPath string,imageFile string,watermarkImageFile string,outputPath string)error {

	err := license.SetMeteredKey("991a1254bef1d5d5ce1729474b8df1ab81e3f94b771e66673aca091193c4b60e")
	 if err != nil {
		 panic(err)
	 }
	
	 // Generate key pair.
	 priv, cert, err := generateKeys()
	 if err != nil {
		 log.Fatal("Fail: %v\n", err)
	 }
 
	 // Create reader.
	 file, err := os.Open(inputPath)
	 if err != nil {
		 log.Fatal("Fail: %v\n", err)
	 }
	 defer file.Close()
 
	 reader, err := model.NewPdfReader(file)
	 if err != nil {
		 log.Fatal("Fail: %v\n", err)
	 }
 
	 // Create the image
	 imgFile, err := os.Open(imageFile)
	 if err != nil {
		 log.Fatalf("Fail: %v\n", err)
	 }
	 defer imgFile.Close()
 
	 signatureImage, _, err := image.Decode(imgFile)
	 if err != nil {
		 log.Fatalf("Fail: %v\n", err)
	 }
 
	 // Create the watermark image
	 wImgFile, err := os.Open(watermarkImageFile)
	 if err != nil {
		 log.Fatalf("Fail: %v\n", err)
	 }
	 defer wImgFile.Close()
 
	//  signatureWImage, _, err := image.Decode(wImgFile)
	//  if err != nil {
	// 	 log.Fatalf("Fail: %v\n", err)
	//  }
	
	 // Create appender.
	 appender, err := model.NewPdfAppender(reader)
	 if err != nil {
		 log.Fatal("Fail: %v\n", err)
	 }
 
	 // Create signature handler.
	 handler, err := sighandler.NewAdobePKCS7Detached(priv, cert)
	 if err != nil {
		 log.Fatal("Fail: %v\n", err)
	 }
 
	 // Create signature.
	 signature := model.NewPdfSignature(handler)
	 signature.SetName("Test Signature Appearance Name")
	 signature.SetReason("Test Signature Appearance Reason")
	 signature.SetDate(time.Now(), "")
 
	 if err := signature.Initialize(); err != nil {
		 log.Fatal("Fail: %v\n", err)
	 }
 
	 numPages, err := reader.GetNumPages()
	 if err != nil {
		 log.Fatal("Fail: %v\n", err)
	 }
 
	 // Create signature fields and add them on each page of the PDF file.
	//  for i := 0; i < numPages; i++ {
		 pageNum := numPages - 1
 
		 // Only Image Signature
		 var opts = annotator.NewSignatureFieldOpts()
		 opts.Rect = []float64{10,10,170,100}
		 opts.Image = signatureImage
		//  opts.WatermarkImage = signatureWImage
 
		 var sigField, err1 = annotator.NewSignatureField(signature, nil, opts)
 
		 sigField.T = core.MakeString(fmt.Sprintf("Signature5 %d", pageNum))
 
		 if err1 = appender.Sign(pageNum, sigField); err1!= nil {
			 log.Fatalf("Fail: %v\n", err1)
		 }
	//  }
 
	 // Write output PDF file.
	 err = appender.WriteToFile(outputPath)
	 if err != nil {
		 log.Fatalf("Fail: %v\n", err)
	 }
 
	 log.Printf("PDF file successfully signed. Output path: %s\n", outputPath)
	 return nil
 }
 
 func generateKeys() (*rsa.PrivateKey, *x509.Certificate, error) {
	 // Generate private key.
	 priv, err := rsa.GenerateKey(rand.Reader, 2048)
	 if err != nil {
		 return nil, nil, err
	 }
 
	 // Initialize X509 certificate template.
	 template := x509.Certificate{
		 SerialNumber: big.NewInt(1),
		 Subject: pkix.Name{
			 Organization: []string{"Test Company"},
		 },
		 NotBefore: now.Add(-time.Hour),
		 NotAfter:  now.Add(time.Hour * 24 * 365),
 
		 KeyUsage:              x509.KeyUsageDigitalSignature,
		 ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		 BasicConstraintsValid: true,
	 }
 
	 // Generate X509 certificate.
	 certData, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	 if err != nil {
		 return nil, nil, err
	 }
 
	 cert, err := x509.ParseCertificate(certData)
	 if err != nil {
		 return nil, nil, err
	 }
 
	 return priv, cert, nil
 }
package models

// FirmaXMLModel representa la firma digital del DTE
type FirmaXMLModel struct {
	XMLName        struct{}      `xml:"Signature"`
	SignedInfo     SignedInfoXML `xml:"SignedInfo"`
	SignatureValue string        `xml:"SignatureValue"`
	KeyInfo        KeyInfoXML    `xml:"KeyInfo"`
}

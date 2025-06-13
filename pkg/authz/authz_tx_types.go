package authz

type TxType string

const (
	InboundVoter  TxType = "InboundVoter"
	OutboundVoter TxType = "OutboundVoter"
	NonceVoter    TxType = "NonceVoter"
	GasPriceVoter TxType = "GasPriceVoter"
)

func (t TxType) String() string {
	return string(t)
}

type KeyType string

const (
	TssSignerKey         KeyType = "tss_signer"
	ValidatorGranteeKey  KeyType = "validator_grantee"
	MuseClientGranteeKey KeyType = "museclient_grantee"
)

func GetAllKeyTypes() []KeyType {
	return []KeyType{ValidatorGranteeKey, MuseClientGranteeKey, TssSignerKey}
}

func (k KeyType) String() string {
	return string(k)
}

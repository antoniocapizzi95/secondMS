package repository

var AddressBookDb AddressBookRepo

func InitModule(abrDb *AddressBookRepo) {
	AddressBookDb = *abrDb
}

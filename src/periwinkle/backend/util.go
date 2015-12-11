// Copyright 2015 Luke Shumaker
// Copyright 2015 Davis Webb

package backend

import (
	"crypto/rand"
	"locale"
	"math/big"
	"periwinkle"
	"strings"

	"github.com/go-sql-driver/mysql"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type errorList []locale.Error

var _ locale.Error = errorList{}

func (e errorList) Locales() []locale.Spec {
	m := map[locale.Spec]int{}
	for _, err := range e {
		for _, l := range err.Locales() {
			m[l] = m[l] + 1
		}
	}
	var ret []locale.Spec
	for l, c := range m {
		if c == len(e) {
			ret = append(ret, l)
		}
	}
	return ret
}

func (errs errorList) L10NString(l locale.Spec) string {
	strs := make([]string, len(errs))
	for i, err := range errs {
		strs[i] = " - " + strings.Replace(err.L10NString(l), "\n", "\n   ", -1)
	}
	return strings.Join(strs, "\n")
}

func (errs errorList) Error() string {
	return errs.L10NString("C")
}

func errHelper(errs *errorList, err locale.Error) {
	if err != nil {
		*errs = append(*errs, err)
	}
}

type table interface {
	dbSchema(*periwinkle.Tx) locale.Error
}

type tableSeed interface {
	table
	dbSeed(*periwinkle.Tx) locale.Error
}

func DbSchema(db *periwinkle.Tx) locale.Error {
	errs := errorList{}
	for _, table := range tables {
		errHelper(&errs, table.dbSchema(db))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func DbDrop(db *periwinkle.Tx) locale.Error {
	errs := errorList{}
	for i := range tables {
		table := tables[len(tables)-1-i]
		errHelper(&errs, locale.UntranslatedError(db.DropTable(table).Error))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func DbSeed(db *periwinkle.Tx) locale.Error {
	errs := errorList{}
	for _, table := range tables {
		if seeder, ok := table.(tableSeed); ok {
			errHelper(&errs, seeder.dbSeed(db))
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Panic, but a little nicer
func dbError(err error) {
	// TODO: return better messages for Conflict errors.
	switch e := err.(type) {
	case sqlite3.Error:
		if e.Code == sqlite3.ErrConstraint {
			panic(periwinkle.Conflict{locale.UntranslatedError(e)})
		}
		panic(locale.UntranslatedError(e))
	case *mysql.MySQLError:
		// TODO: this list of error numbers might not be
		// complete, or totally correct.  See
		// https://mariadb.com/kb/en/mariadb/mariadb-error-codes/
		switch e.Number {
		case 1022, 1062, 1169, 1216, 1217, 1451, 1452, 1557, 1761, 1762, 1834:
			panic(periwinkle.Conflict{locale.UntranslatedError(e)})
		}
		panic(locale.UntranslatedError(e))
	default:
		panic(locale.Errorf("Programmer Error: the programmer said this is a database error, but it's not: %s", e))
	}
}

// Panic, but a little nicer
func programmerError(str string) {
	panic(locale.Errorf("Programmer Error: %s", locale.Sprintf(str)))
}

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var alphabetLen = big.NewInt(int64(len(alphabet)))

func randomString(size int) string {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bigint, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			panic(err) // Luke says this is OK
		}
		bytes[i] = alphabet[bigint.Int64()]
	}
	return string(bytes[:])
}

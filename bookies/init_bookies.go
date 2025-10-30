package bookies

import (
    "diago/utils"
    "reflect"
)

func init() {
    // Use the Bookie interface from utils
    bookieInterface := reflect.TypeOf((*utils.Bookie)(nil)).Elem()

    types := []utils.Bookie{
        &BetGr8{},
        &BetWay{},
        &DimbaKenya{},
        &InBetKenya{},
        &LigiBet{},
        &PariMatch{},
        &SaharaGames{},
        &SportyBet{},
        // Add new bookies here
    }

    for _, b := range types {
        typ := reflect.TypeOf(b)
        if typ.Implements(bookieInterface) {
            utils.Register(b)
        }
    }
}


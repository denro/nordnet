package api

import . "github.com/denro/nordnet/util/models"

type SystemStatusResp SystemStatus

type AccountsResp []Account

type AccountResp AccountInfo

type AccountLedgersResp []LedgerInformation

type AccountOrdersResp []Order

type OrderResp OrderReply

type AccountPositionsResp []Position

type AccountTradesResp []Trade

type CountriesResp []Country

type IndicatorsResp []Indicator

type InstrumentsResp []Instrument

type InstrumentLeverageFilterResp LeverageFilter

type InstrumentOptionPairsResp []OptionPair

type InstrumentOptionPairFiltersResp OptionPairFilter

type InstrumentSectorsResp []Sector

type InstrumentTypesResp []InstrumentType

type ListsResp []List

type LoginResp Login

type LogoutResp LoggedInStatus

type TouchResp LoggedInStatus

type MarketsResp []Market

type NewsResp []NewsPreview

type NewsItemsResp []NewsItem

type NewsSourcesResp []NewsSource

type RealtimeAccessResp []RealtimeAccess

type TickSizesResp []TicksizeTable

type TradableInfoResp []TradableInfo

type TradableIntradayResp []IntradayGraph

type TradableTradesResp []PublicTrades

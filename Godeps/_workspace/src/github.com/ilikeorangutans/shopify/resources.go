package shopify

type Resources interface {
	Apps() *APIPermissions
	Assets(*Theme) *Assets
	Metafields() *Metafields
	Orders() *Orders
	Themes() *Themes
	Transactions() *Transactions
	Webhooks() *Webhooks
}

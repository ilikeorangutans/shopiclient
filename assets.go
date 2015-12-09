package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/ilikeorangutans/shopify"
	"log"
	"os"
	"text/tabwriter"
)

var theme *shopify.Theme
var themeID int

func AssetCommands() cli.Command {
	return cli.Command{
		Name:   "assets",
		Usage:  "Retrieve and modify assets",
		Before: loadThemeForAssets,
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:        "theme-id",
				Usage:       "ID of the theme to use",
				Destination: &themeID,
			},
		},
		Subcommands: []cli.Command{
			cli.Command{
				Name:   "list",
				Usage:  "Lists the assets in a given theme.",
				Action: ListAssets,
			},
			cli.Command{
				Name:   "upload",
				Usage:  "Uploads an asset",
				Action: UploadAsset,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "key",
						Usage: "Key of the asset to upload/modify",
					},
					cli.StringFlag{
						Name:  "value",
						Usage: "Value for the asset",
					},
				},
			},
		},
	}
}

func loadThemeForAssets(context *cli.Context) error {
	themeID := context.Int("theme-id")
	if themeID <= 0 {
		return fmt.Errorf("Invalid theme ID \"%d\"", themeID)
	}

	var err error
	theme, err = shopifyClient.Themes().Get(int64(themeID))
	if err != nil {
		return err
	}

	return nil
}

func ListAssets(context *cli.Context) {
	assets := shopifyClient.Assets(theme)

	assetList, err := assets.List()
	if err != nil {
		log.Fatal(err)
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	writer.Write([]byte("Key\tContent Type\tSize\tUpdated At\n"))
	for i := range assetList {
		asset := assetList[i]
		writer.Write([]byte(fmt.Sprintf("%s\t%s\t%d\t%v\n", asset.Key, asset.ContentType, asset.Size, asset.UpdatedAt)))
	}

	writer.Flush()
}

func UploadAsset(context *cli.Context) {
	key := context.String("key")
	if len(key) == 0 {
		log.Fatal("No key given")
	}

	var asset *shopify.Asset
	value := context.String("value")
	if len(value) > 0 {
		asset, _ = shopify.NewAssetWithValue(key, value)
	}

	result, err := shopifyClient.Assets(theme).Upload(asset)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)
}

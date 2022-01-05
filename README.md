# Flow-Koinly-Export

The provided bash script will convert FT transactions from a single account into the [SIMPLE format](https://help.koinly.io/en/articles/3662999-how-to-create-a-custom-csv-file-with-your-data)

Licensed under MIT.

If you like this script and want to support me send me some FT over at https://find.xyz/bjartek

## Requires
 - curl
 - jq


## How to use
Download the file https://raw.githubusercontent.com/bjartek/flow-koinly-export/main/koinly-simple-export
Execute it like so
```
koinly-simple-export <your_account|.find name>
```

Will result in a csv file in your current directory with the name <your_account>.csv


## Credits
Thanks to https://flowscan.org for a graphql api and permission to use it for this script


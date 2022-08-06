# Flow-Koinly-Export

NB! This does not work anymore after the introduction of flowgraph. Will update this to use flowgraph in the future. 

The provided bash script will convert FT transactions from a single account into the [SIMPLE format](https://help.koinly.io/en/articles/3662999-how-to-create-a-custom-csv-file-with-your-data)

Licensed under MIT.

If you like this script and want to support me send me some FT over at https://find.xyz/bjartek

## Requires
 - curl
 - jq


## Disclaimer
This tool will only export your current transfers in a flow wallet to a csv file that can be imported into koinly on a best effort basis. It is up to the user of this tool to ensure that you follow the tax rules in your country and correct mistakes. 

For instance, for me when I import my transactions all my escrowed versus bids turned up as income and when i get them back they do not turn up as a loss. That is not correct for my countries rules.


## How to use
Download the file https://raw.githubusercontent.com/bjartek/flow-koinly-export/main/koinly-simple-export
Execute it like so
```
koinly-simple-export <your_account|.find name>
```

Will result in a csv file in your current directory with the name <your_account>.csv


Note that if you had any flow before the public launch at 28.01.2020 it will not show in koinly. I have used 0.38 USD per Flow at that time since that is what the ICO auction ended at. 


## Credits
Thanks to https://flowscan.org for a graphql api and permission to use it for this script

## referall

Use my coinly referall link to save some cash https://koinly.io/?via=D607662E

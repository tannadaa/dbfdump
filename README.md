### dbfdump

This is a small utility to dump the data in a directory of DBF files to json.

#### Example

Assume that C:\DBF contains the following files:
* customers.dbf
* invoices.dbf

If you run dbfdump as follows:

`dbfdump -dir C:\DBF`

A new file called dbfdump.json will be created. It will contain data like so:

```json
{
  "customers": [
    {
      "id": 1,
      "name": "customer1"
    }
  ],
  "invoices": [
    {
      "id": 2,
      "amount": 100
    }
  ]
}
```

Each DBF file name becomes a key in the json and each record is an object in the array that key contains.

You can specify the directory for the output file using the `-out` parameter.
If you don't, then the "dbfdump.json" will be created in the current directory.

You can ignore DBF files you don't want included in the json output. Use the `-skip` option.
For example, to ignore the invoices.dbf file from our earlier example, you could use `-skip invoices`.
If you want to skip multiple files, separate them with a comma, but no space, such as `-skip invoices,another,more`.

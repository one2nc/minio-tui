# minio-tui

----
Minio-tui is a terminal-based user interface that allows you to manage your buckets and perform various operations directly from the command line. With its wide range of commands, Minio-tui can save your time by eliminating the need to visit the Minio console to perform these tasks.

----

## To start using minio-tui

```
git clone https://github.com/one2nc/minio-tui
cd minio-tui
go run .
```

## Description

- The Minio-tui dashboard provides an overview of all your buckets, including the total number of buckets you have.
- You can refresh the bucket list by pressing the `r` key.
- To create a new bucket, use the key ` c `.   - - Navigation through the Minio-tui interface is easy, simply use the up and down arrow keys to move between options and press ` enter ` to select an option.
- You can get the presigned url of an object by pressing ` ctrl+p ` and download an object using ` ctrl+d `, that will save the file to the ` ./resources/images ` folder.


`minio-tui`

<img src="./resources/images/minio-tui.png" width="500px" height="" alt="minio-tui">

`Create Bucket` : Press `c`

<img src="./resources/images/createBucket.png" width="500px" height="" alt="minio-tui">

`SearchBucket` : Press `/`

<img src="./resources/images/searchBucket.png" width="500px" height="" alt="minio-tui">

`View Obejct ` : Press `enter` on any bucket

<img src="./resources/images/page2.png" width="500px" height="" alt="minio-tui">

`Presigned URL` : Press `ctrl+p`

<img src="./resources/images/presignedurl.png" width="500px" height="" alt="minio-tui">

`Download Object` : Press `ctrl+d`

<img src="./resources/images/downloadObject.png" width="500px" height="" alt="minio-tui">

# Cloud IP to DB
This project creates a SQLite database containing CIDR records from large Cloud Providers platforms. Currently, the project imports IP CIDR networks from the following:

- Amazon Web Services (AWS)
- Google Cloud (GCP)
- Digital Ocean (DO)
- Microsoft Azure (Azure)
- CloudFlare
- Oracle Cloud

# Technology Stack

The two core elements of this project are:
 - A binary written in GoLang which creates a SQLite database object and populates it with IP CIDR data from various Cloud platform providers.
 - A SQLite database file output containing the Cloud platform providers CIDR information.

The SQLite database schema is made up of a single 'net' table

```CREATE TABLE IF NOT EXISTS net (
 	net_id INTEGER PRIMARY KEY,
 	net TEXT NOT NULL,
 	start_ip INT NOT NULL,
 	end_ip INT NOT NULL,
 	url TEXT NOT NULL,
 	cloudplatform TEXT NOT NULL,
 	iptype TEXT NOT NULL
 	);
 ```
# Use the database for local querying

These instructions describe how to use the sqlite file on your workstation for local querying

## Method 1
The easiest method, simply download the latest cloudIP.sqlite3.db file from the releases section.

https://github.com/stclaird/cloudIPtoDB/releases/latest

Then from the command prompt run the sqlite command and specify the newly download file path.

```
sqlite3 cloudIP.sqlite3.db
```

## Method 2
Use this version to build the binary and create the sqllite database file with the latest cloud provider CIDR information.  This method is slightly more complex but produces the most up to date database file.

Check out this repository
```
git clone git@github.com:stclaird/cloudIPtoDB.git
```
Change directory to the cmd folder
```
cd cloudIPtoDB/cmd/main
```
Build the binary
```
go build -o cloudiptodb main.go
```

Run the binary
```
chmod +x cloudiptodb
./cloudiptodb
```

This will produce a database file in the directory output
You can run this simply with the following command

```
sqlite3 output/cloudIP.sqlite3.db
```

# Querying the Database
The following section gives examples on how to query the database using SQL statements.

### 1. To get the total number of CIDR records across all cloud providers held in the database:

```
select count(*) from net;
```

### 2. To get the number of CIDR records held in the database belong to the cloud platform AWS.

```
select count(*) from net where cloudplatform='aws';
```

or from cloudflare
```
cloudplatformsqlite> select * from net where cloudplatform='cloudflare';
```
```
sqlite> select * from net where cloudplatform='cloudflare';
1729491968|103.21.244.0/22|1729491968|1729492991|https://www.cloudflare.com/ips-v4|cloudflare|IPv4
<SNIP>
sqlite> 
```

3. Find if a specific IP address exists in one of the cidrs held in the database.

The records in the database are in CIDR network format, and not unpacked into individual IP addresses. 
Unpacking the CIDRs into a individual IP records generally doesn't make sense as it would be time consuming to unpack them all and create a much larger than necessary database.
The downside of not having individual IP addresses stored as seperate records will make querying for IPs using SQL difficult.

To remedy this and allow for the querying of individual IP addresses we also store some additonal columns along side the CIDR record, the start (network) and end (broadcast) address. These records are both stored as integers and this means we are able to query whether a specific IP address record exists in the database by testing if the IP address falls between the start record and the end record.

One thing though, for this to work you do need to convert your IP address to an integer before running a query. 
For example, if you want to know if the IP address `177.71.207.129` is within one of the CIDR records stored in the database:

Firstly, you need to convert this IPv4 address to a decimal integer, which is 2974273409. Once you have the IP as an integer you can then perform the following query, where the value for the start_ip and end_ip columns are this integer:

```
SELECT cloudplatform, net 
FROM net 
WHERE start_ip <= '2974273409'
AND end_ip >= '2974273409';
```
If this IP address is contained within one of the CIDR records, this query will return a CIDR record similar to the following,

```
aws|177.71.128.0/17
aws|177.71.207.128/26
```
otherwise, if the IP address is not stored then the database will return no records.



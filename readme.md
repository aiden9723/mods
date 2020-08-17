# Smollan Call Cost Calculator (Callculator)

### Points of Contact

- Aiden Thomas, Internal (can give you more information on this specific request)
- Mpho Jan Vick, External (who sends the requests each month)
- Norma Fundira, Internal (can get you updated rates if required)

### Project Overview

This solution assists our clients in being able to get an estimated cost of an outbound call based on the destination number's network. Once a month (around 14-17), we receive a mail from Mpho Jan Vick from Smollan who requires reporting for the following 'teams' in the Smollan group (and the report's locations on Smollan's XCALLY server):

- Tiger Brands Longmeadow (Located in XCALLY Admin: Analytics > Reports > Custom Reports > Automated Reports > Outbound Calls by Agent - Tiger Brands LM)
- Transnova (Located in XCALLY Admin: Analytics > Reports > Custom Reports > Automated Reports > Outbound Calls - Transnova)
- CIC Headoffice (Located in XCALLY Admin: Analytics > Reports > Custom Reports > Automated Reports > Outbound Call Billing Automated-copy (ID 385))

There are a few steps involved when fulfilling these requests, the details of which will be indicated in the high-level overview, and they are as follows:

1. Generate all 3 reports on XCALLY (Start date 14th of previous month, end date 15th of current month)
2. Download the reports
3. Process reports through application
4. Convert .csv to .xlsx
5. Send to Mpho

### System Architecture (or what's involved)

- XCALLY Extracted Custom Reports - Used to generate & download the extracted reports
- Call cost calculator application - Written in Golang, application that runs in the command line

### High Level Overview

The application runs in the commandline, it accepts 2 inputs:

`report`, which is a downloaded extracted .csv report from XCALLY

`rates`, which is a .csv file that contains a list of number codes, the network they belong to, and what it costs to make a call to that network per minute

Once given accepted `rates` and `report` files, the application will output a .csv file, with the `report` filename and `-formatted` appended. The `report` file needs to be a .csv and have the following headers in the first row, in this exact order:

1. `UniqueID` - Unique ID
2. `Destination` - Destination Number
3. `Billsec` - Billable Seconds
4. `Tag` - Call Tag
5. `Prefix` - Call Prefix
6. `Agent Name` - Agent who made the call
7. `Context` - Telephony-specific Call Context
8. `Date` - Date the call was made

Thus far, the `rates` pricing hasn't been updated since implementation, and so the headers accepted by the application needs to be a .csv with the following order of columns, without a row for headers:

1. Number code
2. Country of Origin
3. Network Name
4. Cost p/min, in South African Rands



### Code Base Organisation

> Word of warning: This isn't a very robust (or easily readable) codebase, as the time from planning to release was very short, and on short notice. Unless the requirement changes, or if you're looking for a challenge to try improve your development skills by reading badly written and implemented code, don't bother too much with continuously improving this repo. 

The application goes through several steps, taking the initial report, and creating a new .csv report file which contains each call detail record from the `report` input file, with the following headers appended to the first row:

`Code` - The network associated to the call destination number (which is checked by means of RegEx)

`Cost` - The ESTIMATED cost of the call, by calculating the associated cost of the originated call to the destination's network 

In a nutshell, the following sequence is followed when the application has accepted both input files:

1. Read `report` file
2. Append headers `Code` and `Cost` to the output .csv file
3. Read `rates ` file
4. Separate valid/billable `report` call destination numbers from invalid/non-billable destination numbers
5. Filter through calls, discarding calls with 0 billable seconds
6. Match/validate numbers from `report` against all codes in `rates` file
7. Calculate estimated cost of call with the matching network cost of the destination number
8. Add newly calculated cost of the call and associated network to the record and its own list
9. Write modified records to new output file and append `-formatted` to the filename, preserving its original file type (.csv)



In the `main.go` file, you'll find more detail into the underlying logic, the packages used and their underlying logic.



### Version Control Procedures

None, the code base is hosted on the `scopserv-southafrica` GitHub account, feel free to implement.

### QA Workflow

Nothing outside of testing done during development.




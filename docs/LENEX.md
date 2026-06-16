---
title: "Lenex 3.0 - Technical documentation"
subtitle: "An international data exchange format for swimming"
version: "3.0"
author: "Christian Kaufmann, Switzerland"
contact: "lenex@swimrankings.net"
source: "https://www.swimrankings.net/lenex"
copyright: "© 1999-2025 - All Rights Reserved"
document_date: "2025-03-04"
pdf_creation_date: "2025-03-04T20:51:36+01:00"
pdf_mod_date: "2025-03-04T20:51:36+01:00"
pdf_creator: "Microsoft® Word für Microsoft 365"
pdf_producer: "Microsoft® Word für Microsoft 365"
pdf_pages: 35
source_file: "Lenex_3.0_Technical_Documentation.pdf"
transcribed: true
---

# Lenex 3.0 - Technical documentation

An international data exchange format for swimming

## 1. General

A "Lenex" file is a XML file with some additional constraints, which cannot be defined with an XSD schema.
Because of this, and to give programmers an easier to read document with additional information on the
format, we have put together this "non-standard" kind of documentation for you.

Chapter 9 at the end of the document lists all changes done since release of Lenex 3.0 and contains some
details about the history of Lenex.

## 2. Updates and versions

Lenex has remained more or less unchanged since version 3.0. Compatibility is important when exchanging
data. Therefore, almost all changes since 2008 are extensions only. There are very few minor breaking
changes.

Some names no longer correspond to the terms used today (e.g. FINA, LEN, handicap), but are kept like this
in order to avoid breaking changes.

## 3. Lenex files

A Lenex file is a XML file with the extension `.lef`. Usually, Lenex files are compressed (in the ZIP file format)
and are labeled with the extension `.lxf`. A Lenex file can contain all kinds of data at the same time. However,
when exchanging data, the following files with a subset of possible objects are commonly used:

- **Invitation**: An invitation file contains general information, the schedule and the event structure of one
  meet. Additionally, it could be necessary or helpful to add time standards and/or qualification times, which
  are important for the meet.
- **Entries**: An entry file contains the entries for one meet. One file might contain the entries of one club only,
  or it can contain all entries of all clubs.
- **Startlist**: Almost the same as an entry file. Normally it contains all entries of all clubs, but it is possible to
  split the seeded entries of each session into separate files. To define the seeding, all ENTRY elements
  need the attributes `heatid` and `lane`.
- **Results**: A result file contains the results of one meet. Normally, it contains all results of all clubs, but it is
  possible to split the results for each club into a separate file.
- **Records**: A record file contains one or more list(s) of records.
- **Time standards**: A time standards file may contain different kinds of time standards and/or qualification
  times. It might make sense to store time standards in separate files, if they are independent of meets (e.g.
  Olympic A and B time standards). If the time standards are bound to a certain meet, they should be
  included in the invitation file of that meet.

## 4. Lenex tree

This is an overview of the Lenex structure. Subchapter 4.1 is a tree with the most important elements. The
other subchapters describe some of the sub trees in more details. To get the full information about all
elements that are allowed or required, please refer to the chapter with the "Lenex element documentation"
(chapter 5).

The items in a Lenex tree fall into three categories:

- **Elements**: An element contains any number of child objects and attributes.
- **Collections**: A collection is an element that contains elements of one type only. By default, the collection
  will be given the plural name of the element name it contains (e.g. SESSIONS contains SESSION
  objects). In most cases, the objects contained in a collection have at least one required attribute for
  identification purposes (e.g. attribute "distance" for element SPLIT).
- **Attributes**: Attributes contain data in one of the basic Lenex data formats. The recognized formats are
  documented in the chapter "Lenex data types" (chapter 6). Attributes can be attached to objects or to
  collections.

### 4.1. Tree overview

The following tree shows the most important elements in a Lenex tree.

```
LENEX                                             The root of a Lenex file.
  CONSTRUCTOR                                     Information about the creator of the file.
  MEETS
     MEET                                         The root for a meet sub tree.
       SESSIONS                                   The schedule and event details of a meet.
       CLUBS
          CLUB                                    All data of one club at the meet.
            ATHLETES
            COACHES
            RELAYS
            OFFICIALS
  RECORDLISTS                                     The root for the record lists sub tree.
  TIMESTANDARDLISTS                               The root for the time standard lists sub tree.
```

### 4.2. Sub tree SESSIONS

The SESSIONS sub tree describes the entire event structure with prelims and final events and the age
groups used for the result lists.

```
SESSIONS
  SESSION                                         Data of one session with all its events.
     POOL
     EVENTS
       EVENT                                      Description of one event/round.
          AGEGROUPS
             AGEGROUP                             Details of one age group.
                RANKINGS
                  RANKING                         Details for ranking with reference to result elements.
          HEATS
             HEAT                                 Details of one heat (number, starttime).
          SWIMSTYLE
          TIMESTANDARDREFS
             TIMESTANDARDREF                      Reference to a timestandard list.
     JUDGES
       JUDGE                                      Details about judges for a session.
```

### 4.3. Sub tree ATHLETES

The ATHLETES sub tree contains all athletes of one club with their entries and/or results.

```
ATHLETES
  ATHLETE                                        Data of one athlete.
    ENTRIES
       ENTRY                                     Entry for one event/round.
         MEETINFO
    RESULTS
       RESULT                                    Result for one event/round.
         SPLITS
```

### 4.4. Sub tree RELAYS

The RELAYS part is used to describe relay entries and results of one club. The relay swimmers are not
stored directly in this tree. A unique id is stored in the tree in order to reference an athlete in the ATHLETES
sub tree.

```
RELAYS
  RELAY                                          Data of one relay team.
    ENTRIES
       ENTRY                                     Entry for one event/round.
         RELAYPOSITIONS
           RELAYPOSITION
              MEETINFO
    RESULTS
       RESULT                                    Result for one event/round.
         RELAYPOSITIONS
           RELAYPOSITION
         SPLITS
```

### 4.5. Sub tree RECORDLISTS

The sub tree RECORDLISTS is used to define all kind of records. One record list contains all records of a
specific type (e.g. world records), gender and pool length (course). In this sub tree, the information about the
athletes is represented by means of complete ATHLETE objects, and not just as a reference to some other
sub tree.

```
RECORDLISTS
  RECORDLIST                                     Data of one record list (type, gender, course).
    AGEGROUP
    RECORDS
      RECORD                                     Data of one record (individual or relay).
         SWIMSTYLE
         ATHLETE
         RELAY
           RELAYPOSITIONS
              RELAYPOSITION
                ATHLETE
         MEETINFO
         SPLITS
```

### 4.6. Sub tree TIMESTANDARDLISTS

The sub tree TIMESTANDARDLISTS is used to define time standards and qualification times.

```
TIMESTANDARDLISTS
   TIMESTANDARDLIST                                 Time standards (type, gender, course).
      TIMESTANDARDS
         TIMESTANDARD                               Data of one time standard / qual. time.
            SWIMSTYLE
```

## 5. Lenex element documentation

The following list is alphabetically ordered and describes the meaning and content of every element in a
Lenex tree. Element names are in uppercase. Attribute names are in lowercase. The basic data types are
described in chapter "Lenex data types" (chapter 6). Elements can appear in different ways:

- **Normal (`-`)**: Zero or one instance of an attribute/element is allowed.
- **Required (`r`)**: Exactly one instance of an attribute/element is required.
- **Multiple (`m`)**: There can be any number of instances of an elements (including zero). An attribute is
  always allowed once in maximum per element.

Every element or collection can have an attribute with the name `[elementname]id` (e.g. `athleteid` for the
element ATHLETE). For some objects, this attribute is mandatory, because it is used to build relationships
between elements in different sub trees. Attributes have to be unique over all instances of the same element
type.

---

### Element AGEDATE

The AGEDATE is the date used to calculate the age of an athlete.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `type` | e | r | The type describes, how the age is calculated. The following values are acceptable: `YEAR` (the age is calculated using the year of the meet and the year of birth only), `DATE` (the age is calculated exactly between the date and the birth date), `POR` (age calculation according the Portuguese federation), `CAN.FNQ` (calculation according the Quebec federation), `LUX` (calculation according the Luxembourg federation). |
| `value` | d | r | The date value. |

### Element AGEGROUP

This element contains information about an age range. It is used in events and record lists.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `agegroupid` | n | r | Only for events, every AGEGROUP element needs an id, because the objects can be referenced from ENTRY objects. The id has to be unique within an AGEGROUPS collection. |
| `agemax` | n | r | The upper bound of the age range. -1 means no upper bound. |
| `agemin` | n | r | The lower bound of the age range. -1 means no lower bound. |
| `gender` | e | - | In mixed events, the gender can be specified in the AGEGROUP objects. Values can be male (M), female (F) or mixed (X). Mixed is only used for relays. This can be useful to define events with gender set to all (A), but the ranking is separated. This attribute is not allowed in the context of a RECORDLIST or TIMESTANDARDLIST element. |
| `calculate` | e | - | Information for relay events about how the age is calculated: `SINGLE` (this is the default value; the age of each relay swimmer has to be in the given range), `TOTAL` (the total age of all swimmers has to be in the given range). |
| `handicap` | e | - | The sport class for the agegroup. This is used to group results by disability categories. Allowed values are: 1 - 15, 20, 34, 49 standard sport classes. |
| `levelmax` | s | - | The maximum level (A-Z) of the agegroup. If the value is missing, this means no upper bound. |
| `levelmin` | s | - | The minimum level (A-Z) of the agegroup. If the value is missing, this means no lower bound. |
| `levels` | s | - | comma separated list of codes of allowed athlete levels. |
| `name` | s | - | The name of the age group (e.g. "Juniors"). |
| RANKINGS | o | - | A collection with references to results ranked in this agegroup. |

### Collection AGEGROUPS

This collection contains all age group definitions of one event.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| AGEGROUP | o | m | The definition of one age group. |

### Element ATHLETE

This contains all information of a athlete including all entries and results in the context of a meet sub tree.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `athleteid` | n | r | The id attribute should be unique over all athletes of a meet. It is required for ATHLETE objects in a meet sub tree. |
| `birthdate` | d | r | The date of birth for the athlete. If only the year of birth is known, the date should be set to January 1st of that year. |
| CLUB | o | - | The club or team for the athlete, when he swam the record. This element is allowed in a recordlist sub tree only. |
| ENTRIES | o | - | All entries of the athlete. This element is allowed in a meet sub tree only. |
| `firstname` | s | r | The first name of the athlete. |
| `firstname.en` | si | - | The first name in english. |
| `gender` | e | r | Gender of the athlete. Values can be male (M), female (F) or non binary (N). |
| HANDICAP | o | - | Information about the sport classes of a swimmer. |
| `lastname` | s | r | The last name of the athlete. |
| `lastname.en` | si | - | The last name in english. |
| `level` | s | - | The level of the athlete (used with levels in AGEGROUP). |
| `license` | s | - | The registration number given by the national federation. This number should be looked at together with the nation of the club the athlete is listed in the Lenex file. |
| `license_ipc` | n | - | The registration number given by World Para Swimming, also known as SDMS ID. |
| `nameprefix` | s | - | An optional name prefix. For example for Peter van den Hoogenband, this could be "van den". |
| `nation` | e | - | See table "Nation Codes" for acceptable values. |
| `passport` | s | - | The passport number of the athlete. |
| RESULTS | o | - | All results of the athlete. This element is allowed in a meet sub tree only. |
| `status` | e | - | The following values are acceptable: `EXHIBITION` (the athlete swims exhibition in all events), `FOREIGNER` (the athlete is a foreigner), `ROOKIE`. |
| `swrid` | n | - | The global unique athlete id given by swimrankings.net. |

> Federation extensions (see §7): Germany adds `license_dbs` (German para swimming federation id) and `license_dsv` (German federation id, DSV). Switzerland adds status value `SUI.STARTSUISSE` (special registration for foreign athletes allowed to swim in a relay for a Swiss Record).

### Collection ATHLETES

This collection contains all athletes of one club.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| ATHLETE | o | m | Data of one athlete. |

### Element BANK

This is used to represent bank payment information for a meet.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `accountholder` | s | r | The name of the bank account holder. |
| `bic` | s | - | THE BIC code of the bank. |
| `iban` | s | r | The IBAN number of the bank account. Must be a valid IBAN number. |
| `name` | s | - | The name of the bank. |
| `note` | s | - | Note for the payment as additional information. |

### Element CLUB

In the meet sub tree, this element contains information about a club, including athletes and relays with their
entries and/or results. In the record list sub tree, the element contains information about the club or nation of
record holders.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| ATHLETES (1) | o | - | The athletes of this club. |
| `code` | s | - | The official club code given by the national federation. Only official club codes should be used here! |
| CONTACT (1) | o | - | Contact address for the specific meet. |
| `name` | s | r | The full name of the club or the team. |
| `name.en` | si | - | The club name in english. |
| `nation` | e | - | See table "Nation Codes" for acceptable values. |
| `number` (1) | n | - | This value can be used to distinguish different teams of the same club in a meet entries or results file. |
| OFFICIALS (1) | o | - | The officials from this club. |
| `region` | s | - | The code of the regional or local swimming committee. Only official codes should be used here! |
| RELAYS (1) | o | - | The relay teams of this club. |
| `shortname` | s | - | A short version of the club name. This string is limited to 20 characters. |
| `shortname.en` | si | - | The short version of the club name in english. |
| `swrid` | n | - | The global unique club id given by swimrankings.net. |
| `type` | e | - | The following types of clubs are allowed: `CLUB` (this is the default value), `NATIONALTEAM` (the club represents a national team of a federation; in this case, the code, region and nation attribute should be the same), `REGIONALTEAM` (the club represents a regional team; in this case, the code and region attribute should be the same), `UNATTACHED` (to be used for the CLUB entry, that contains data of athletes, where the club is unknown; in this case, the attribute name and the CONTACT element are not required). |

(1) These objects and elements are not used in CLUB objects, which appear in the record list sub tree.

### Collection CLUBS

This collection contains all clubs that take part of one meet.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| CLUB | o | m | The data of one club. |

### Element COACH

This element contains all information about an coach.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| CONTACT | o | - | The contact information of the coach. |
| `firstname` | s | r | The first name if the coach. |
| `gender` | e | - | Gender of the official. Values can be male (M) female (F) or non binary (N). |
| `lastname` | s | r | The last name of the coach. |
| `nameprefix` | s | - | An optional name prefix. For example for Peter van den Hoogenband, this could be "van den". |
| `nation` | e | - | See table "Nation Codes" for acceptable values. |
| `passport` | s | - | The passport number of the coach. |
| `type` | e | - | The following values are acceptable: `HEADOFMISSION` (for "chef de mission"), `HEADCOACH` (the head coach of the club), `COACH` (any other coach), `STAFF` (any other staff member of the club). |

### Collection COACHES

This collection contains all coaches and staff of a club.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| COACH | o | m | The data of one coach. |

### Element CONSTRUCTOR

This element contains information about the software, which created the Lenex file and the contact
information about the provider of that software.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| CONTACT | o | r | Contact information of the provider of the software, which created the Lenex file. |
| `name` | s | r | Name of the application that created the Lenex file. |
| `registration` | s | - | Name of user, to who the creator application was registered. |
| `version` | s | r | The version number of the application that created the Lenex file. |

### Element CONTACT

This element contains the contact address for a person or organisation.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `city` | s | - | The city of the contact address. |
| `country` | s | - | See table "Country codes" for acceptable values. |
| `email` | s | r | The e-mail address of the contact. The attribute is required in the context of a CONSTRUCTOR element only. |
| `fax` | s | - | The fax number of the contact. |
| `internet` | s | - | The full URL of the website of the contact person or organisation. The https:// should be included in the string. |
| `name` (1) | s | - | The full name of the contact person or the name of the organisation. |
| `mobile` | s | - | The mobile phone number of the contact person. |
| `phone` | s | - | The phone number of the contact person or the organisation. |
| `state` | s | - | The state, province or county of the contact address. |
| `street` | s | - | The first additional line of the address. |
| `street2` | s | - | The second additional line of the address. |
| `zip` | s | - | The postal code of the address. |

(1) These elements are not used in the context of an OFFICIAL element.

### Collection ENTRIES

This collection contains all entries of on athlete or a relay team.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| ENTRY | o | m | The data of one entry. |

### Element ENTRY

This element contains the information for a single entry of an athlete or a relay to a specific round of a meet.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `agegroupid` | n | - | Reference to an age group (AGEGROUP element in the AGEGROUPS collection of the EVENT element). |
| `entrycourse` | e | - | This attribute indicates a pool length for the entry time. This is necessary when special seeding rules are used. See section 6.4 for acceptable values. |
| `entrydistance` | n | - | The entry distance in centimeters. Is used for some fin swimming events. For such entries the entrytime should be "NT". |
| `entrytime` | st | - | The entry time in the swim time format. |
| `eventid` | n | r | Reference to the EVENT element using the id attribute. |
| `handicap` | e | - | In special cases, the sport class can be different for a single entry. Allowed values are the values for standard sport classes. |
| `heatid` | n | - | Reference to a heat (HEAT element in HEATS collection of the EVENT element). |
| `lane` | n | - | The lane number of the entry. |
| MEETINFO | o | - | This element contains the information, about a qualification result for the entry time was achieved. |
| RELAYPOSITIONS | o | - | Only for relay entries. This element contains references to the relay swimmers. |
| `status` | e | - | This attribute is used for the entry status information. An empty status attribute means a regular entry. The following values are allowed: `EXH` (exhibition swim), `RJC` (rejected entry), `SICK` (athlete is sick), `WDR` (athlete/relay was withdrawn). |

The combination of the attributes `eventid`, `heatid` and `lane` should be unique over all ENTRY objects of the
same meet.

### Element EVENT

This element contains all information of an event. For events with finals, there has to be an EVENT element
for each round.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| AGEGROUPS | o | - | The AGEGROUPS collection contains the descriptions for the age groups in this event. For Open/Senior events, AGEGROUPS is only needed with one AGEGROUP element as a placeholder for the RANKINGS element (for places in result lists). If `round="FHT"`, then no AGEGROUPS element is allowed. |
| `daytime` | t | - | The daytime of the start of the event. |
| `eventid` | n | r | Every event needs to have an id attribute, so that it can be referenced by ENTRY and RESULT objects. The id attribute has to be unique over all EVENT objects of all sessions of a meet. |
| FEE | o | - | The entry fee for this event. If there are global fees per athlete, relay and/or meet, the FEE elements in the MEET element should be used. |
| `gender` | e | - | The gender of the event. The default value is all (A). Other values allowed are male (M), female (F) and mixed (X). Mixed is only used for relays. |
| HEATS | o | - | Collection with all heats in the event. |
| `maxentries` | n | - | The maximum number of entries per club in this event. To limit the number of entries per athlete or relay, use the maxentries attribute in the MEET element. |
| `number` | n | r | The number of the event. The event numbers should be unique over all events of a meet. The EVENT objects of the different rounds for the same event may have the same event number. |
| `order` | n | - | This value can be used to define the order of the events within a session if it is not by the event number and if there are no start times for the events. |
| `preveventid` | n | - | This value is a reference to a previous event's id (e.g. the prelims events for final events). The default value is -1 and means, that there was no previous event. |
| `round` | e | - | The following values are allowed here: `TIM` (this is the default value; used for an event with timed finals), `FHT` (fastest heats of an event with timed finals; events with this value for round should always refer to the corresponding timed final event, which should be of the same distance, stroke and age groups; events with round set to FHT only make sense for the schedule and ENTRY objects, but never to be used for RESULTs), `FIN` (used for finals including A, B, C, … finals), `SEM` (for semi finals), `QUA` (for quarterfinals), `PRE` (for prelims), `SOP` (swim-off after prelims), `SOS` (swim-off after semi-finals), `SOQ` (swim-off after quarterfinals), `TIMETRIAL` (time trial event). |
| `run` | n | - | Used if there is more than one swim-off necessary. Default value is 1. |
| `status` | e | - | Used for incremental data files. Possible values are: `ENTRIES` (entries are available for this event), `SEEDED` (event is seeded), `RUNNING` (event is the current running event), `UNOFFICIAL` (some heats have results but are not all official yet), `OFFICIAL` (all heats of the event are done and official). |
| SWIMSTYLE | o | r | The SWIMSTYLE element contains information about distance and stroke of the event. |
| TIMESTANDARDREFS | o | - | A list of references to TIMESTANDARDREF elements with references to time standard lists to be used for this event. |
| `timing` | e | - | The type of timing for an event. If missing, the session should be checked and finally the value for the meet should be used. See MEET for acceptable values. |
| `type` | e | - | The following values are allowed: the default value is empty (this applies for regular events that are run according to World Aquatics rules), `MASTERS` (master results to be used for master rankings/records). |

> Federation extension (Germany, §7): `round` has additional value `GER.RES` ("Nachschwimmen"/re-swim, used for German team championships).

### Collection EVENTS

This collection contains all events of one session.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| EVENT | o | m | The data of one event. |

### Element FACILITY

This element contains name and full address of meets facility.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `city` | s | r | The city of the facility. |
| `nation` | s | r | See table "Country codes" for acceptable values. |
| `name` | s | - | The name of the facility (e.g. "Aquatic Center"). |
| `state` | s | - | The state, province or county of the address. |
| `street` | s | - | The first additional line of the address. |
| `street2` | s | - | The second additional line of the address. |
| `zip` | s | - | The postal code of the address. |

### Element FEE

The fee is used in MEET and EVENT objects.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `currency` | e | - | See table "Currency Codes" for acceptable values. |
| `type` (1) | e | r | Used in context of FEES in MEET or SESSION objects only. Acceptable values are: `CLUB` (global fee to be paid per club for the meet), `ATHLETE` (global fee to be paid per athlete), `RELAY` (global fee to be paid per relay team), `TEAM` (global fee to be paid per team, e.g. Swiss Team Championship), `LATEENTRY.INDIVIDUAL` (fee per entry for individual late entries; used in the context of MEET.FEES only), `LATEENTRY.RELAY` (fee per entry for relay late entries; used in the context of MEET.FEES only). |
| `value` | c | r | The value of the fee in the currency format. |

(1) This element is required only in the context of a FEES collection.

### Collection FEES

This collection contains all global fees in a MEET element.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| FEE | o | m | The data of one fee. |

### Element HANDICAP

This element is used to describe sports classes for an athlete.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `breast` | e | r | The sport class for breaststroke. Allowed values are: 0 - 15 standard sport classes, `GER.AB` (swimmers with a minor disability), `GER.GB` (swimmers with a mental disability). |
| `breaststatus` | e | - | The state of the sport class for breaststroke. Allowed values are: `NATIONAL` (the sport class is on national level only), `NEW` (the sport class is not yet valid), `REVIEW` (the sport class must be reviewed in this year, but is valid up to the end of the year), `OBSERVATION` (the sport class needs the observation during the meet in order to get the confirmed status), `CONFIRMED` (the sport class is confirmed for international meets). |
| `exception` | s | - | Additional information beside the sport classes. |
| `free` | e | r | The sport class for freestyle, backstroke and fly. Allowed values are the same as for the breast attribute. |
| `freestatus` | e | - | The state of the sport class for freestyle, backstroke and fly. Same values used as for breaststatus. |
| `medley` | e | r | The sport class for individual medley. Allowed values are the same as for the breast attribute. |
| `medleystatus` | e | - | The state of the sport class for individual medley. Same values used as for breaststatus. |

### Element HEAT

The heat is used to define more details in the start list (e.g. schedule).

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `agegroupid` | n | - | Reference to an age group (AGEGROUP element in the AGEGROUPS collection of the EVENT element). |
| `daytime` | t | - | The daytime of the start of the event. |
| `final` | e | - | This value is used to identify A, B, ... finals. Allowed are characters A, B, C and D. |
| `heatid` | n | r | The id attribute should be unique over all heats of a meet. It is required when you have ENTRY / RESULT objects that reference to a heat. |
| `number` | n | r | The number of the heat. The heat numbers have to be unique within one event, also in a case when you have A finals in different agegroups. |
| `order` | n | - | This value can be used to define the order of the heats within an event if it is not by the heat number and if there are no start times for the heats. |
| `status` | e | - | The status of the heat. Allowed values are: `SCHEDULED` (the heat is scheduled but not seeded yet), `SEEDED` (the heat is seeded), `RUNNING` (heat is the current running heat), `UNOFFICIAL` (results are available but not official yet), `OFFICIAL` (results of the heat are official). |

### Collection HEATS

This collection contains all heats of one event.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| HEAT | o | m | The data of one heat. |

### Element JUDGE

This element contains information to attach an official to a session with his role in the session.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `number` | n | - | The number for judges where there are more than one. This can be used for example for the lane number for timekeepers. |
| `officialid` | n | r | A reference to a OFFICIAL element. |
| `remarks` | s | - | Additional information for the judge. |
| `role` | e | - | Indicates the role of a judge. The list is built according to World Aquatics descriptions. See table "Judge roles" for acceptable values. |
| `status` | e | - | Indicates the status of a judge. Acceptable values are: `AVAILABLE` (for entries to provide the information about where an official may be entered with what function), `SCHEDULED` (final information how the planning for judges looks like), `COMPLETED` (in the results how the jury finally was in place). |

### Collection JUDGES

This collection contains all judges of one session.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| JUDGE | o | m | The data of one judge. |

### Element LENEX

This is the root element of every Lenex file which identifies it as a XML file conforming to the Lenex data
format.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| CONSTRUCTOR | o | r | This element contains information about the software which created the Lenex file. |
| `created` | ts | - | The timestamp when the file was created (UTC time). |
| MEETS | o | - | Contains all the information of meets like athletes, relays, entries and results. |
| RECORDLISTS | o | - | Contains different types of records (e.g. World records, Olympic records) including age group records. |
| `revisiondate` | d | - | This is a reference to the revision date of this document. If the value is missing, the file references to the latest version of the Lenex format. |
| TIMESTANDARDLISTS | o | - | Contains different type of time standards and qualification times. |
| `version` | s | r | The version number of the Lenex format. The current value is "3.0". |

### Element MEET

This element contains all information of one meet, including events, athletes, relays, entries and results.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| AGEDATE | o | - | The date to be used to calculate the age of athletes. The default value is the date of the first session and type by year of birth only. |
| BANK | o | - | Information for a bank account for payment of entry fees. |
| `altitude` | n | - | Height above sea level of the meet city. |
| `city` | s | r | The name of the city where the meet was run. Should be the same as FACILITY.city. |
| `city.en` | si | - | Name of meet city in english. |
| CLUBS | o | - | Collection of clubs of the meet. |
| CONTACT | o | - | The contact address of the meet organizer. |
| `course` | e | - | The size of the pool. See section 6.4 for acceptable values. If the attribute is not available, all SESSION objects need to have a course attribute. |
| `deadline` | d | - | The date for the entry deadline. |
| `deadlinetime` | t | - | The exact time for the entry deadline. |
| `entrystartdate` | d | - | The date from when (online) entries are open/available. |
| `entrytype` | e | - | The type of (online) entries: `OPEN` (the meet is open for all clubs), `INVITATION` (the meet is open to invited clubs only). |
| FACILITY | o | - | The full address of the meets facility. |
| FEES | o | - | Fees used for this meet. On this level, different global fees for clubs, athletes and relays are allowed. If there are fees that have to be paid per entry, the FEE element in the EVENT objects should be used. |
| `hostclub` | s | - | The executing federation or club of the meet (e.g. the German Swimming Federation, if the European Champ was held in Berlin). |
| `hostclub.url` | s | - | A website url, that directs to the executing club for the meet. |
| `maxentriesathlete` | n | - | The maximum number of individual entries per athlete in this session. |
| `maxentriesrelay` | n | - | The maximum number of relay entries per club in this session. |
| `name` | s | r | The name of the meet. Normally the name should not contain a full date (maybe the year only) and/or a city or pool name. |
| `name.en` | si | - | Meet name in english. |
| `nation` | e | r | The three letter code of the nation of the meet city. This should be the same as FACILITY.nation. |
| `number` | s | - | The sanction number for the meet by the federation. |
| `organizer` | s | - | The organisation which promotes the meet (e.g. World Aquatics for the World Championships). |
| `organizer.url` | s | - | A website url, that directs to the organizer of the meet. |
| POINTTABLE | o | - | Description of the point table used for scoring. |
| POOL | o | - | Details about the pool where the meet took place. |
| QUALIFY | o | - | Details about how qualification times for entries are defined. |
| `reservecount` | n | - | The number of reserve swimmers in finals and semifinals. |
| `result.url` | s | - | A website url, that directs to the results page. This should be a deep (direct) link to the result lists and not the general website of a meet. |
| SESSIONS | o | r | Description of all events grouped by session. |
| `startmethod` | e | - | The following values are possible: `1` (one start allowed), `2` (two starts allowed). |
| `status` | e | - | Used for incremental data files. Possible values are: `ENTRIES` (entries are available for meet), `SEEDED` (at least one event is seeded), `RUNNING` (at least for one event results are available), `OFFICIAL` (all events are done and official). |
| `swrid` | uid | - | The global unique meet id given by swimrankings.net. |
| `timing` | e | - | The type of timing for a meet. Acceptable values are: `AUTOMATIC` (a full automatic timing system was used), `SEMIAUTOMATIC` (automatic start, manual button for end time), `MANUAL3` (timing was done with three manual times per lane), `MANUAL2` (timing was done with two manual times per lane), `MANUAL1` (timing was done with one manual time per lane). |
| `touchpadmode` | e | - | Information about timing installation for a meet. Acceptable values are: `ONESIDE` (touchpads on one side of the pool), `BOTHSIDE` (touchpads on both sides of the pool). |
| `type` | e | - | The meet type. The following values are allowed: the default value is empty (this applies for normal meets that are run according to World Aquatics rules), all other values depend on definitions of national federations. |
| `withdrawuntil` | d | - | The date for withdrawals from the entry list. |

### Element MEETINFO

This element is used in entries and records for general information about a meet.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `approved` (2) | s | - | Contains a code for the organisation, who approved the qualification time, e.g. World Aquatics, European Aquatics or a IOC nation code. If this field is empty, the qualification time was not approved. |
| `city` (1) | s | r | The city name where the meet took place. |
| `course` (2) | e | - | This attribute indicates the pool length, where the qualification time was achieved. See section 6.4 for acceptable values. |
| `date` (1) | d | r | The date of the swim of the record or qualification time achievement. |
| `daytime` | t | - | The day time of the swim. |
| `name` | s | - | The meet name. |
| `nation` (1) | e | r | The nation of the city for the meet. |
| POOL | o | - | The details about the pool. |
| `qualificationtime` (2) | st | - | The qualification time, since this can be different to the entry time. If the value is missing, the entry time is the qualification time. |
| `state` | s | - | The state of the city for the meet. |
| `timing` | e | - | The type of timing. See MEET for acceptable values. |

(1) These elements are required only in the context of a RECORD element.
(2) These elements are used in the context of a ENTRY / RELAYPOSITION element only.

### Collection MEETS

This collection allows you to put the results of more than one meet in the same Lenex file. However, our
experience with Lenex during the last years shows that it is better to keep different meets in separate files.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| MEET | o | m | Data of one meet. |

### Element OFFICIAL

This element contains all information about an official.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| CONTACT | o | - | The contact information of the official. |
| `firstname` | s | r | The first name if the official. |
| `gender` | e | - | Gender of the official. Values can be male (M) and female (F). |
| `grade` | s | - | The grade of the judge. The value here is specific to national federations and depends on their officials education system. |
| `lastname` | s | r | The last name of the official. |
| `license` | s | - | The registration number given by the national federation. |
| `nameprefix` | s | - | An optional name prefix. For example for Peter van den Hoogenband, this could be "van den". |
| `nation` | e | - | See table "Nation Codes" for acceptable values. |
| `officialid` | n | r | The id attribute should be unique over all officials of a meet. It is required for JUDGE objects in a meet sub tree. |
| `passport` | s | - | The passport number of the official. |

### Collection OFFICIALS

This collection contains all officials of a club.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| OFFICIAL | o | m | The data of one official. |

### Element POINTTABLE

This element is used to describe the point scoring used for a meet.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `name` | s | r | The name of the point score system. |
| `pointtableid` | n | - | Common point tables have a unique id. Details are in chapter 6.5 ("Ids for point scores"). |
| `version` | s | r | The version number/year of the point score system. |

### Element POOL

This element is used to describe the pool where the meet took place.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `lanemax` | n | - | Number of the last lane used in the pool for the meet. The number of lanes can be calculated with LANEMAX - LANEMIN + 1. |
| `lanemin` | n | - | Number of the first lane used in the pool for the meet. |
| `temperature` | n | - | The water temperature. |
| `type` | e | - | The type of the pool. Acceptable values are INDOOR, OUTDOOR, LAKE and OCEAN. |

### Element QUALIFY

This element contains information about details, how qualification entrytimes are defined.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `conversion` | e | - | The way, how times are converted for seeding: `NONE` (this is the default value; no conversion will be done), `FINA_POINTS` (entry times, flagged with a course other than the event course will be converted by calculating the World Aquatics points and from there calculating back to a swim time), `PERCENT_LINEAR` (here the conversion will be done by adding / subtracting a certain percentage to the original entry time), `NON_CONFORMING_LAST` (in this case, entries that are flagged with the event course will be seeded first, all other entries will be seeded after these). |
| `from` | d | r | The first day of the qualification period for entry times. |
| `percent` | n | - | The percentage for conversion PERCENT_LINEAR. |
| `until` | d | - | The last day of the qualification period for entry times. If it is missing, then the default is the last day before the first day of the meet. |

### Element RANKING

This element describes one entry in the rankings of one agegroup. It contains the place and a reference to a
result (individual or relay).

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `medal` | e | - | In some cases the value from attribute place with the round given in EVENT is not accurate enough to determine if medals were awarded or not. For this you can use this attribute to indicate, if there was a medal awarded in the ranking for a specific result. The following values are allowed here: `GOLD`, `SILVER`, `BRONZE`, `NONE` (no medal was awarded, e.g. in the ranking of a prelim event). |
| `order` | n | - | This value can be used to define the order of the results. If it is missing, the value for place is used to sort the elements in a collection. |
| `place` | n | r | The final position in the result list for the current event/round. |
| `resultid` | n | r | A reference to the RESULT element. |

### Collection RANKINGS

This collection contains a set of ranking elements.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| RANKING | o | m | Data of one item in ranking. |

### Element RECORD

This element describes one individual or relay record. It is possible to have no ATHLETE / RELAY objects. In
this case the record is a "record standard time".

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| ATHLETE | o | - | The person who holds the record. This is only used for individual records. |
| `comment` | s | - | This value can be used for additional comments like "Swum in the prelims" or "Ratification pending at World Aquatics". |
| MEETINFO | o | - | Information about the meet, where the record was swum. |
| RELAY | o | - | The relay team and swimmers, who holds the record. This is only used for relay records. |
| SPLITS | o | - | The split times of the record. |
| SWIMSTYLE | o | r | The swimstyle contains information like distance, stroke of the record. |
| `swimtime` | st | r | The final time of the record in the swim time format. |
| `status` | e | - | State of the record. The following values are allowed: `APPROVED` (record is approved and valid), `PENDING` (waiting for ratification by federation), `INVALID` (invalid due ratification failures), `APPROVED.HISTORY` (valid, but not current record anymore), `PENDING.HISTORY` (waiting, but not current record anymore), `TARGETTIME` (no record yet but a time to define the minimal time for a new record). |

> FAQ: The ATHLETE / RELAY element is not required in a RECORD. None of both elements should be present when there is no record but only a required (target) time.

### Element RECORDLIST

This element describes one single record list.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| AGEGROUP | o | - | For agegroup records. Agegroup is "Open", if the element is missing. |
| `course` | e | r | The course for the record list. See section 6.4 for acceptable values. |
| `gender` | e | r | The gender for records in this list. Acceptable values are male (M), female (F) and mixed (X). Mixed is only used for relays. |
| `handicap` | e | - | The sport class for the record list. Allowed values are: 1 - 15, 20, 34, 49 standard sport classes. |
| `name` | s | r | The name of the record list (e.g. "World Records"). |
| `nation` | s | - | For international records, this field is empty. For national or regional records, the field should contain the IOC three letter code of the national federation. |
| `order` | n | - | This value can be used to define an order for all recordlists within a collection. |
| RECORDS | o | r | The records of this record list. |
| `region` | s | - | For international and national records, this field is empty. For regional records, the field should contain the official code for the region. If region has a value, nation needs to have a value as well. |
| `updated` | d | - | The date of the last change to the record list. |
| `type` | e | - | The record type. The following values are allowed: `WR` (world records), `OR` (Olympic records), `ER` (European records), `PAR` (Pan American records), `AFR` (African records), `AR` (Asian records), `OCR` (Oceanian records), `CWR` (Commonwealth records), IOC three letter nation code (the national records of the specific federation). This list of types may be extended by federations. In this case, values should have the nation code and a dot as prefix (e.g. `SUI.RZW` for records of a region in Switzerland). |

### Collection RECORDLISTS

This collection contains a set of record lists. For each different combination of gender, course, age group or
type, a separate RECORDLIST element is needed.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| RECORDLIST | o | m | Data of one record list. |

### Collection RECORDS

This collection contains all records of a record list. If there is more than one athlete holding the same record,
each of them has a RECORD element in the collection.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| RECORD | o | m | Data of one individual or relay record. |

### Element RELAY

This element is used to describe one relay team for a record or for a meet.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `agemax` (2) | n | r | The maximum age allowed for the oldest swimmer in the team. The value -1 means no upper bound. |
| `agemin` (2) | n | r | The minimal age allowed for the youngest swimmer in the team. The value -1 means no lower bound. |
| `agetotalmax` (2) | n | r | The maximum total age of all swimmers in the relay team. The value -1 means that the total age is unknown. |
| `agetotalmin` (2) | n | r | The minimum total age of all swimmers in the relay team. The value -1 means that the total age is unknown. |
| CLUB (1) | o | - | The club or team of the relay in the context of a record. |
| ENTRIES (2) | o | - | All entries of the relay team. |
| `gender` (2) | e | r | The gender of the relay team. Acceptable values are male (M), female (F) or mixed (X). |
| `handicap` (2) | e | - | For relays with swimmers with disability. The default value is 0. Other values allowed are: `14` (relay with athletes with a mental disability), `20` (relay athletes of a total of 0-20 points from sport classes), `34` (relay athletes of a total of 21-34 points from sport classes), `49` (relay athletes of a total of 35 points or more or blind athletes). |
| `name` | s | - | The name of the relay team. |
| `number` (2) | n | - | The team number of the relay team. This number is only necessary, if there are different teams of the same club in the same agegroups / events. |
| RELAYPOSITIONS (1) | o | - | The relay swimmers in the context of a relay record. |
| RESULTS (2) | o | - | All results of the relay team. |

(1) These objects are allowed in the context of a record only.
(2) These elements/objects are allowed in the context of a meet only.

### Collection RELAYS

This collection contains all relays of one club of a meet.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| RELAY | o | m | Data of one relay team. |

### Element RELAYPOSITION

This element is used for information about one relay swimmer.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| ATHLETE | o | - | Last name, first name, etc. of the athlete. This element is allowed in the context of a record only and in this case it is required. |
| `athleteid` | n | - | A reference to the ATHLETE element of the athlete. This attribute is allowed in the context of a meet sub tree only. |
| MEETINFO | o | - | This element contains the information, where the entry time was achieved. This element is only allowed in the context of a relay entry. |
| `number` | n | r | The number of the swimmer in the relay. The first swimmer is 1, the second 2 and so on. -1 can be used to add reserve swimmers. |
| `reactiontime` | rt | - | The reaction time at the start of the first swimmer and the relay take over times for other swimmers. |
| `status` | e | - | No status attribute means the swimmer finished his part correctly. Otherwise, the following values are allowed: `DSQ` (relay athlete was disqualified), `DNF` (relay athlete did not finish). |

### Collection RELAYPOSITIONS

This collection contains information's about relay swimmers.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| RELAYPOSITION | o | m | Data of one relay swimmer. |

### Element RESULT

This element is used to describe one result of a swimmer or relay team.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `comment` | s | - | Additional comment e.g. for new records or reasons for disqualifications. |
| `eventid` | n | r | Reference to the EVENT element using the id attribute. |
| `handicap` | e | - | In special cases, the sport class can be different for a single result. Allowed values are the values for standard sport classes. |
| `heatid` | n | - | Reference to a heat (HEAT element in HEATS collection of the EVENT element). |
| `lane` | n | - | The lane number of the entry. |
| `points` | n | - | The number of points for the result according to the scoring table used in a meet. |
| `reactiontime` | rt | - | The reaction time at the start. For relay events it is the reaction time of the first swimmer. |
| RELAYPOSITIONS | o | - | The information about relay swimmers in this result. Only allowed for relay RESULT objects. |
| `resultid` | n | r | Each result needs a unique id which should be unique over a meet. |
| `status` | e | - | This attribute is used for the result status information. An empty status attribute means a regular result. The following values are allowed: `EXH` (exhibition swim), `DSQ` (athlete/relay disqualified), `DNS` (athlete/relay did not start — no reason given or too late withdrawal), `DNF` (athlete/relay did not finish), `SICK` (athlete is sick), `WDR` (athlete/relay was withdrawn on time). |
| SPLITS | o | - | The split times for the result. In a Lenex file, split times are always saved continuously. |
| `swimdistance` | n | - | The result distance in centimeters. Is used for some fin swimming events. For such results the swimtime should be "NT". |
| `swimtime` | st | r | The final time of the result in the swim time format. |

### Collection RESULTS

This collection contains all results of a athlete or relay team of a meet.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| RESULT | o | m | Data of one single result. |

### Element SESSION

This element is used to describe one session of a meet.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `course` | e | - | With indicating a pool length per session, the global value of the meet can be overridden, e.g. if the prelim sessions are short course and the finals are long course. See section 6.4 for acceptable values. |
| `date` | d | r | The date of the session. |
| `daytime` | t | - | The daytime when the session starts. |
| `endtime` | t | - | The time when the session ended. |
| EVENTS | o | r | The events of the session. |
| FEES | o | - | Fees used for this session. On this level, different global fees for clubs, athletes and relays are allowed. If there are fees that have to be paid per entry, the FEE element in the EVENT objects should be used. |
| JUDGES | o | - | The judges of the session. |
| `maxentriesathlete` | n | - | The maximum number of individual entries per athlete in this session. |
| `maxentriesrelay` | n | - | The maximum number of relay entries per club in this session. |
| `name` | s | - | Additional name for the session e.g. "Day 1 - Prelims". |
| `number` | n | r | The number of the session. Session numbers in a meet have to be unique. |
| `officialmeeting` | t | - | The daytime when the officials meeting starts. |
| POOL | o | - | The details about the pool, if they are different per session. Otherwise use the element in MEET. |
| `remarksjudge` | s | - | Additional remarks given by the referee. |
| `status` | e | - | Used for incremental data files. Possible values are: `ENTRIES` (entries are available for this session), `SEEDED` (at least one event of the session is seeded), `RUNNING` (it is the session with the current running event), `UNOFFICIAL` (some events have results but are not all official yet), `OFFICIAL` (all of the session events are done and official). |
| `teamleadermeeting` | t | - | The daytime when the team leaders meeting starts. |
| `timing` | e | - | The type of timing for a session. If missing, the global value for the meet should be used. See MEET for acceptable values. |
| `touchpadmode` | e | - | Information about timing installation for the session. See MEET for acceptable values. If not available, value from element MEET is used. |
| `warmupfrom` | t | - | The daytime when the warmup starts. |
| `warmupuntil` | t | - | The daytime when the warmup ends. |

### Collection SESSIONS

Depending on the context this collection contains all sessions of a meet or all sessions for a judge, where he
is planned in.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| SESSION | o | m | Data of one session. |

### Element SPLIT

This element contains information about a single split time. In a Lenex file, split times are always saved
continuously.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `distance` | n | r | The distance where the split time was measured. |
| `swimtime` | st | r | The time of the result in the swim time format. |

### Collection SPLITS

This collection contains all available split times for a single result.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| SPLIT | o | m | Data of one split time. |

### Element SWIMSTYLE

This element is used to describe one swim style.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `code` | s | - | A code (max. 6 characters) of the swim style if the stroke is unknown. |
| `distance` | n | r | The distance for the event. For relay events it is the distance for one single athlete. |
| `name` | s | - | The full descriptive name of the swim style if the stroke is unknown (e.g. "5 x 75m Breast with one Arm only"). |
| `relaycount` | n | r | The number of swimmers per entry / result. Value 1 means, that it is an individual event. All other values mean, that it is a relay event. |
| `stroke` | e | r | The following values are allowed here: `BACK` (backstroke), `BREAST` (breaststroke), `FLY` (fly or butterfly), `FREE` (freestyle), `IMRELAY` (for relays where each athlete swims all strokes like in an individual medley event), `MEDLEY` (for individual and relay medley; the order of stroke is according to World Aquatics rules: Fly, back, breast, free for individual events; back, breast, fly, free for relay events), `UNKNOWN` (for all special events; in this case, the name attribute of the event is mandatory). For fin swimming there are the following values: `APNEA`, `BIFINS` (bi-fins), `DYNAMIC` (result is in meter, attribute swimdistance), `DYNAMIC_BIFINS` (dynamic with bi-fins, result is in meter, attribute swimdistance), `DYNAMIC_NOFINS` (dynamic without fins, result is in meter, attribute swimdistance), `IMMERSION`, `MULTIPLE` (for relays, according to CMAS rules: bi-fins, bi-fins, surface, surface), `SPEED_APNEA`, `SPEED_ENDURANCE`, `STATIC`, `SURFACE`. |
| `swimstyleid` | n | - | The id attribute is important for SWIMSTYLE objects, where the stroke attribute is "UNKNOWN". In this case, the id should be a unique value to help to identify the swim style. |
| `technique` | e | - | The technique of the style. If this attribute is missing or empty, it means normal swimming. All other values are mainly used for technical events in meets for kids. The following values are allowed: `DIVE` (for swimming under water), `GLIDE` (for gliding only), `KICK` (for kick only — swimming with legs, no arms used), `PULL` (for pull only — swimming with arms, no legs used), `START` (for start only), `TURN` (for turn only). |

> Federation extension (Germany, §7): `stroke` has additional value `GER.APH` ("Apnoe mit Hebeboje").

### Element TIMESTANDARD

This element describes one time standard.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| SWIMSTYLE | o | r | The style contains information like distance, stroke of the record. For each TIMESTANDARD in the same collection, the SWIMSTYLE should be unique. |
| `swimtime` | st | r | The time standard or qualification time. |

### Collection TIMESTANDARDS

This collection contains a set of time standards.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| TIMESTANDARD | o | m | Data of one time standard (time and swim style). |

### Element TIMESTANDARDLIST

This element describes one single time standard list.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| AGEGROUP | o | - | For age group time standards. Agegroup is "Open", if the element is missing. |
| `course` | e | r | The course for the timestandard list. See section 6.4 for acceptable values. |
| `gender` | e | r | The gender for time standards in this list. Acceptable values are male (M), female (F) and mixed (X). Mixed is only used for relays. |
| `handicap` | e | - | The sport class for the timestandard list. Allowed values are: 1 - 15, 20, 34, 49 standard sport classes. |
| `name` | s | r | The name of the time standard list (e.g. "Olympic A Time Standards"). |
| `timestandardlistid` | n | r | The unique id of the time standard list. |
| TIMESTANDARDS | o | r | The time standards or qualification times of this list. |
| `type` | e | - | There can be different type of time standards. Default value is MAXIMUM. The following values are allowed: `DEFAULT` (the time standards describe a set of default times, that may be used in a team competition, where a result for an event of a team is missing or invalid and therefore is replaced by this default time), `MAXIMUM` (the time standards describe a maximal time for a meet; swimmers can only compete, when they are faster than the time standards), `MINIMUM` (the time standards describe a minimal time for a meet; swimmers can only compete, when they are slower than the time standards). |

### Collection TIMESTANDARDLISTS

This collection contains a set of time standard lists. For each different combination of gender, course, age
group or type, a separate TIMESTANDARDLIST element is needed.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| TIMESTANDARDLIST | o | m | Data of one time standard list. |

### Element TIMESTANDARDREF

This element describes a reference from a meet to a time standard list.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| `timestandardlistid` | n | r | The id of the time standard list element. |
| FEE | o | - | An optional element with a fine for missed time standards. |
| `marker` | s | - | An optional string to be used to mark the result, if the time standard was missed. Or to mark a result if a qualification time was fulfilled. |

### Collection TIMESTANDARDREFS

This collection contains a set of time standard references.

| Element/Attribute | Type | Card. | Remarks |
|---|---|---|---|
| TIMESTANDARDREF | o | m | Data of one time standard reference. |

## 6. Lenex data types

In a Lenex file, the following data types are used:

| Type | Code | Description |
|---|---|---|
| Currency | c | An integer number. Currency values are represented in cents, e.g. one dollar in the Lenex currency format is 100. |
| Date | d | Dates are always represented by a string in the form "YYYY-MM-DD". Example: "2004-03-09" means March 9, 2004. |
| Daytime | t | A daytime (hour and minutes) represented by a string in the form "HH:MM". Hours should be from 0 to 24, minutes from 0 to 59. |
| Enumeration | e | An enumeration is a set of predefined values that are allowed in the attribute of that data type. |
| Number | n | A signed 32-bit integer number. Only the characters "0".."9" and "-" are allowed. |
| Global Identifier | uuid | A unique identifier (128bit label). Format should be a string with 36 characters (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx). |
| String | s | A string containing any character. Special characters like `< > " '` and `&` have to be quoted with `&lt; &gt; &quot; &apos;` and `&amp;`. For line ends in multiline strings you should use `&#10;` (encoding for the LF character). |
| String international | si | Same as "String" but only characters between ASCII #32 - #127 allowed. |
| Swim time | st | The swim time data type is always a fixed length string of the following form: "HH:MM:SS.ss". HH: hours from 0 to 99, MM: minutes from 0 to 59, SS: seconds from 0 to 59, ss: hundredths of a second from 0 to 99. Example: "00:14:45.86" means a time of 14:45.86. In addition the string "NT" is allowed if no time is available. |
| Timestamp | ts | A timestamp in ISO 8601 format "YYYY-MM-DDThh:mm.ss". By default, UTC time should be used. |
| Reaction time | rt | All reaction times are numbers and are measured in hundredths of a second. The first character indicates, if the reaction time is positive (+) or negative (-). Example "+14" means a positive reaction time of 14 hundredths. The reaction time "0.00" should be transmitted as "0". If the reactiontime is missing, the value should be empty. |
| Unique id | uid | Unique id's are a character (A-Z) followed by a number. Additional separator characters (space, dash, point) are allowed but have no meaning and should be ignored when comparing unique id's. The id's are unique worldwide and are handled by swimrankings.net. |

### 6.1. Other enumeration definitions

Some definitions for valid values in enumerations are held in a separate Excel sheet:
[https://www.swimrankings.net/files/LenexCodes.xlsx](https://www.swimrankings.net/files/LenexCodes.xlsx)

The following tables exist:

- **Nation codes**: For the nation codes, the three letter codes of IOC are used.
- **Country codes**: The international two letter postal codes are used.
- **Currency codes**: The international three letter codes are used.
- **Judge roles**: A code is defined for each common possible role.
- **Ids for point scores**: For common point scorings there is a list of id's that can / should be used for the POINTTABLE elements.

### 6.4. Course codes

For the `course` attribute, the following values are allowed: `LCM`, `SCM`, `SCY`, `SCM16`, `SCM20`, `SCM33`,
`SCY20`, `SCY27`, `SCY33`, `SCY36` and `OPEN` for open water swimming.

## 7. Specific extensions for different federations

### 7.1. Germany (GER)

Element ATHLETE can have the attribute `license_dbs` for the registration id of the german para swimming
federation. `license_dsv` is used for the german federation id (DSV).

For the `round` attribute in EVENT there is an additional value: `GER.RES`. This means "Nachschwimmen" or
re-swim and is used for the German team championships.

For `stroke` attribute in SWIMSTYLE there is an additional value `GER.APH` for "Apnoe mit Hebeboje".

### 7.2. Switzerland (SUI)

For the `status` attribute in ATHLETE there is an additional value: `SUI.STARTSUISSE`. This is a special
registration for foreign athletes to be allow to swim in a relay for a Swiss Record.

## 8. Frequently asked questions (FAQ)

**Q: What should I put in the required attributes for an ATHLETE in a RECORD when there is no record but a required time only?**

A: The ATHLETE / RELAY element is not required in a RECORD. None of both elements should be there in
such a case.

## 9. History

### 9.1. Changes since introduction version 3.0

| Date | Changes |
|---|---|
| 4 Mar 2025 | BANK.accountholder has to be mandatory (due new SEPA payment rules). ATHLETE: Attribute gender allows N for non binary. COACH / COACHES added in CLUB. JUDGE: New value status added. Moved tables for countries, currencies, nations and roles to a new source (Excel file). |
| 10 Feb 2025 | SWIMSTYLE: New stroke MULTIPLE for fin swimming added. Attribute status added in elements MEET, SESSION and EVENT. New value for attribute status in HEAT. |
| 10 Dec 2024 | LENEX.revisiondate added. RANKING.medal added. |
| 2 Dec 2024 | LENEX.created added. New datatypes added (timestamp, uuid). |
| 27 Nov 2024 | HEAT.status: New value SCHEDULED added. MEETINFO: Attribute timing added. HANDICAP: Attributes breaststatus, freestatus and medleystatus added. ENTRY and RESULT: Attribute handicap added. RECORD.status: Changed from comment to enumeration. |
| 7 Oct 2024 | New element BANK in MEET added. New values for fin swimming in stroke in SWIMSTYLE added. swimdistance in RESULT and entrydistance in ENTRY added (fin swimming). |
| 1 Oct 2024 | license_ipc, license_dbs for element ATHLETE added. |
| 20 May 2024 | Attribute status for element ATHLETE added. |
| 03 Mar 2023 | New round value TIMETRIAL added. |
| 02 Nov 2021 | New value for stroke added: IMRELAY for relays where each athlete swims all strokes. |
| 27 Oct 2021 | Attribute timing (MEET, SESSION, EVENT) with new values SEMIAUTOMATIC and MANUAL2. |
| 10 Nov 2016 | POOL: name removed. MEET: new element FACILITY added. |
| 21 Oct 2016 | MEET: maxentries removed, replaced by maxentriesathlete and maxentriesrelay. |
| 10 Mar 2016 | SWIMSTYLE: BIFINS was missing as value for "stroke". |
| 05 Dec 2015 | MEET: attributes reservecount, startmethod and withdrawuntil added. SESSION: attribute endtime, maxentriesathlete, maxentriesrelay and remarksjudge added. JUDGE: attribute remarks added. Definition for strings extended for multiline strings. |
| 27 Jul 2015 | SESSION: attribute timing added. |
| 31 Jan 2015 | MEET.FEES.FEE: New values for type attribute added. MEET: attribute entrystartdate and entrytype added. |
| 22 May 2013 | SESSION: FEES collection added. |
| 05 Dec 2012 | AGEGROUP: attribute levels added. ATHLETE: attribute level added. |
| 06 Nov 2012 | Attribute "middlename" removed everywhere. |
| 30 Mar 2011 | AGEGROUP: type attribute removed. EVENT: type attribute (for master events) added. TIMESTANDARDLIST: timestandardlistid is mandatory attribute. TIMESTANDARD: swimstyle should be unique in a TIMESTANDARDS collection. |
| 09 Mar 2010 | ENTRY: Attribute entrycourse added. HEAT: Attribute heatid is mandatory as it is referred in ENTRY elements. QUALIFY: removed relaymode, changed conversion, added percent. |
| 05 Oct 2009 | RESULT: Attribute resultid is mandatory as it is referred in RANKING elements. MEET: Attributes result.url, organizer.url, hostclub.url added. |
| 17 Aug 2009 | MEET: organizer revised, promoter removed, hostclub added. |
| 21 May 2009 | LCM_BEFORE_SCM added for attribute conversion in QUALIFICATION. |
| 15 Dec 2008 | QUALIFICATION.relaymode: new value BOTH added. |
| 27 Nov 2008 | HEATS, RANKINGS, TIMESTANDARDREFS added. English attributes added (e.g. "firstname.en"). In AGEGROUP "type" attribute renamed to "calculate", new type attribute added. FEE elements in MEET grouped in FEES collection. heatid and lane added to RESULT element. deadlinetime added in MEET element. In SWIMSTYLE, technique, values DIVE, GLIDE and PULL added. |
| 25 Nov 2008 | Introduction version 3.0. |

### 9.2. History before version 3.0

Lenex 1.0 was introduced in September 1999.

In February 2004 Lenex 2.0 format was introduced. It is based on Lenex 1.0, but there is no direct
compatibility between the two formats. However, the basic structure of events, entries and results data is the
same. Generally the following changes were made:

- Added new sub trees for record lists and time standards.
- Added and refined existing objects for entries and results.
- Some of the key words were renamed in order to have them in a more accurate English and to avoid
  short named unclear identifier.
- Some data types were changed in order to have no country specific separators anymore.
- Some of the key words in enumerations were changed to remove unclear short abbreviations.
- All basic types are specified as attributes now. Entities are used for complex data types only.

Lenex 3.0 is based on format 2.0 and compatibility between the two is more or less ensured. This means an
application, that can handle format 2.0, can read data from format 3.0 with some minor exceptions. Lenex
3.0 covers the following extensions:

- HEATS added in EVENT to reflect a more complete schedule and event structure. An ENTRY / RESULT
  has no heat attribute anymore, but a reference to a HEAT element (heatid).
- RANKINGS added in AGEGROUP to reflect more complete result lists with different overlapping
  rankings. The place and agegroupid attribute in RESULT element are deprecated. Instead the new
  RANKING elements refer to a RESULT element with their resultid attribute. For events with no
  agegroups, you still need to define at least one AGEGROUP to hold the list with the places. Exhibition
  swims are just saved as result for an ATHLETE / RELAY and have no reference in any RANKINGS
  collection.
- For language dependent attributes (e.g. firstname / lastname) it is possible to add optional international
  version of attributes (e.g. firstname.en) with the english text.
- In AGEGROUP renamed "type" to "calculate". New "type" attribute to identify master events.
- The EVENT can have a TIMESTANDARDREFS collection with elements that refer to a
  TIMESTANDARDLIST and the corresponding fine. Like this it can be defined exactly for which events
  timestandards do apply.
- FEE elements in MEET are now grouped in a FEES collection.
- Attribute "swrid" added in MEET, CLUB and ATHLETE for global unique identification of these elements.
- The preferred default encoding should be UTF-8 now. However, ansi codepage encoding is still possible.
  The encoding used for the file should always be set correctly in the prolog.
- Other minor adjustments.

---

*© 1999-2025 - All Rights Reserved — [https://www.swimrankings.net/lenex](https://www.swimrankings.net/lenex)*
*Author: Christian Kaufmann, Switzerland — lenex@swimrankings.net*

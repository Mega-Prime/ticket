# User stories

## MVP

* Create a ticket

As a customer, I go to a web page and I fill in the details of my issue. I click 'Submit'.

* List tickets by status

As an engineer, I want to be able to list all tickets.

* See a ticket by ID

As an engineer (or customer) I can type the ticket ID into a web page and I will see the corresponding ticket.

* Edit a ticket

As a customer I want to see my ticket with any updates, and add a comment, or close the ticket.

As an engineer I need to be able to look at an issue, add information, modify the status, and maybe close the ticket.

## What don't we have for MVP?

* WebApp (1wk)
  * Web page for new issue
  * Web page for list tickets
  * Web page for get ticket by ID
  * Web page for edit issue
  * Persistent storage - done
  * API server to open/save ticket data

## Engineering work for MVP

* Web app
* ListTickets API endpoint (1d)
    * Write Get all tickets test
    * Write Get all Tickets func  

* UpdateTicket API endpoint (2d)


* Database integration (2wk)
  * Create interface for Store - done
  * Refactor existing Store code as e.g. 'MemoryStore' - done
  * Add DB-backed Store implementation: 'MongoStore' - done
    * Add integration test for MongoStore AddTicket - done
      * Convert database Object ID to ticket ID (1d) - done
    * Add integration test for MongoStore GetticketByID (1d)

* Demo Day (2wk) - Oct 6th 2020

## Database

* Translation layer like an API
* Endpoints: create ticket, read a ticket, update a ticket

Different implementations for different DB engines: Mongo, SQLite, Postgres

## Open questions

Are issues private? Can a customer see other customers' issues? Can all engineers see all tickets, or is there any kind of access control? OAuth? User logins?

## Extensions or enhancements

* Numeric ticket IDs (or otherwise a bit more friendly than 12-byte hex strings)

* and see a message like: "Thanks for reporting your issue, and here's a link to check on the status of it. We'll email you when it's resolved."

* filter tickets by status

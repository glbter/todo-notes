# todo-notes
test project of todo notes system

task: to create organizer-calendar service (RESP API) with the functionality to:
 - add events, reminders.
 - modify/update them, change title/name, time, description...
 - remove events
 - view the list of events for a day, week, mounth, year (with filtering by features)
 
2nd step:
 - add user authorization (save credentials locally (in-memory for now), will move them to DB later. Considering all safety measures (do not store plain passwords wink)
 - all endpoints are closed with authorization 
 - user has his own timezone (and an API to change it)
 - all tasks have to be viewed in the user's timezone. (and created considering the user's timezone or provided timezoneони (+1 parameter)
 
3rd step:
 - add graceful shutdown
 - metrics endpoint
 - proper logs (enough information for debug without violation of security)
 - test coverega  
 
4th step:
 - connection to DB (Postgres)
 - database migration (setup proper structure)
 - integration via docker-compose
 - replace storage from in-memory to DB
 - add integration tests to cover logic with DB storage

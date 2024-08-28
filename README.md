# Mapperizer
### Website Map Generator

[mapperizer.codethings.net](https://mapperizer.codethings.net/?demo=true)

Create a map for your website with a google spreadsheet.      


- Create google spreadsheet with specific field headers (or [clone it from here](https://docs.google.com/spreadsheets/d/1KX8kjoQ7RxFDh0T8roQ83JwIyd8w8etRIJfdJljCook/edit?gid=0#gid=0))
- Add your locations   
- Publish and sync your spreadsheet at [mapperizer.codethings.net](https://mapperizer.codethings.net)
- Paste your iframe snippet into your website 


## Stack
- Server API built with Go
- Data cached in Redis
- Client API endpoints called with HTMX
- Styled with TailwindCSS
- Map instance built with Leaflet
- Dockerized and deployed to render.com 
- Custom GPT for quickly sourcing location data 

## Location Grabber GPT
https://chatgpt.com/g/g-kJ2i1rFDl-location-grabber    
Worked for a bit on this custom GPT to source location data.  The GPT accepts a broad set of possible location descriptors, determines specific entities from these descriptors, scrapes the web for location data for each entity and returns structured result data.  
 
The structured response includes a table of location metadata (address, web, phone, lat/long) plus a code snippet that can be copy pasted into the google spreadsheet.

So far surprisingly accurate results.  
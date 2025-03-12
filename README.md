# photo-map
Small project to display marker on the world map, which show a photo gallery of your photos when clicking on them

You can test it [here](https://travel.ax4w.me)

> You are free to fork the project and change / adjust anything to your needs :^) 

> Note: If the container is stopped **while** generating thumbnails, then you need to manually start `generate.sh`

## Setup
1. Pull the docker image
2. Mount a volume onto `/app/map-data` (**inside that folder the folder structure and scripts will be stored where you need to store your images**)
3. Configure the enviroment variables listed under **environment**
4. Start the Container
   
### Environment
You will need to set the following environment variables, which are used to connect to your postgres database
- `dbname` : the name of the database
- `host` : ip of the database server
- `port` : port of the database server
- `user` : username to log into the database server
- `password` : password to log into the database server

  
## Adding Locations
Folders and scripts are automatically created when starting the container.

Simply create a folder inside the `images` folder, which is inside the mounted volume, with the name of the location (e.g `Berlin` , `Munich`,...).

---
### Example
```
/home/map-data (<- the folder the volume is mounted to in my casewin)
-> images/
--> Berlin
---> pic1.jpg
---> pic2.jpg
-> thumbs/
...
```
---

Inside that folder you can place the photos you want to display at that marker. 

After max. 3 Minutes a new entry should be visible on the map. If the marker is at the wrong location you will need to adjust `lat` and `long` manually in the database.
The `lat` and `long` are fetched from [nominatim](https://nominatim.openstreetmap.org/) and could be wrong if there are multiple locations with same name.

The gallery might look broken, because each time images change the thumbnails are being rebuild. If you have many images stored the rebuilding might take longer and the 
thumbnails might not be generated yet. So just wait a bit :-) (You can see in the logs, if generation has been started / is finished)

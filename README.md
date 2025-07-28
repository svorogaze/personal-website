# Portfolio & Blog Website  
A full-stack personal portfolio and blog website built with Next.js (frontend), Go (backend), and deployed with Docker.  
## Tech Stack  
- **Frontend**: Next.js  
- **Backend**: Go  
- **Database**: PostgreSQL  
- **File Storage**: Minio  
- **Reverse Proxy**: Nginx  
- **Deployment**: Docker 
## Deployment
#### Option 1: Using Pre-Built Docker Images 
  1. Copy compose_images.yaml
  2. Create a '.env' file based on 'example.env'
  3. (Optional)Create init_db.sql file
  4. Create folders for minio and nginx
  5. Create nginx config inside nginx folder
  6. Run: ```docker-compose -f compose_images.yaml up```
####  Option 2: Building images locally
  1. Download repo
  2. Create a '.env' file based on 'example.env'
  3. (Optional) Create init_db.sql file
  4. Run: ```docker-compose -f compose.yaml up```

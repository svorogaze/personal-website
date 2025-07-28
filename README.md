A modern personal portfolio and blog website
### Uses:
- Frontend: NextJS
- Backend: go
- Reverse proxy: nginx
- S3 file storage: minio
- Database: PostgreSQL
- Docker
### Deploying:
#### To deploy with images from my dockerhub repo:
  1. Copy compose_images.yaml
  2. Create a '.env' file following 'example.env'
  3. Create init_db.sql file(if needed)
  4. Create folders for minio and nginx
  5. Create nginx config inside nginx folder
  6. Do docker-compose
#### To deploy with images built locally
  1. Download repo
  2. Create a '.env' file following 'example.env'
  3. Create init_db.sql file(if needed)
  4. Do docker-compose

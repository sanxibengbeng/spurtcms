
<p align="center">
  <a href="https://www.spurtcms.com/#gh-light-mode-only">
      <img src="https://www.spurtcms.com/public/img/v2-logo.png" width="318px" alt="Spurtcms logo" />
  </a>
   
</p>
<h3 align="center">Open Source Golang based CMS Solution - Self hosted </h3>
<p align="center"> Build with Golang + PostgreSQL</p>


<br />
<p align="center">
  <a href="https://github.com/spurtcms/spurtcms/releases">
    <img src="https://img.shields.io/github/last-commit/spurtcms/spurtcms/main" alt="GitHub last commit" />
  

  </a>
  <a href="https://github.com/spurtcms/spurtcms/issues">
    <img src="https://img.shields.io/github/issues/spurtcms/spurtcms/main" alt="GitHub issues" />
    
  </a>

  <a href="https://github.com/spurtcms/spurtcms/releases">
    <img src="https://img.shields.io/github/repo-size/spurtcms/spurtcms" alt="GitHub repo size" />
             
  </a>
  


</p>
<br />

> [!IMPORTANT]
> 🎉 <strong>Spurtcms V1.2.0 is now available!</strong> Read more in the <a target="_blank" href="https://www.spurtcms.com/spurtcms-change-log" rel="dofollow"><strong>announcement post</strong></a>.[![GoDoc](https://godoc.org/github.com/spurtcms/block?status.svg)](https://pkg.go.dev/search?q=spurtcms)
<br />
<p>
spurtCMS Admin Panel prioritizes user-friendly administration, offering powerful tools for content creation, management, and defining CMS workspaces. Administrators have precise control over member access, ensuring streamlined member management. Dynamic channel management allows effective content structuring, enhancing the overall user experience. Administrators effortlessly create and manage channels and templates, providing a comprehensive, user-centric content management solution for personalized and organized web environments.
</p>
<br />

![GIF](https://dev.spurtcms.com/public/img/animated-gif-maker.gif)


## ❯  🚀 Easy to Deploy Spurtcms Admin Panel on your server

This is the official repository of Spurtcms.You can easily deploy spurtcms in your local server.

### Step 1: Download the source files:

Clone the Git repository that contains spurtCMS Admin project files, and .env file from the path https://github.com/spurtcms/spurtcms using the “git clone” command.

```
git clone https://github.com/spurtcms/spurtcms
```
After successful git clone, you should see a folder “spurtcms-admin” with folders locales, view, storage, public and  .env, file.


### Step 2: Database Setup

spurtCMS supports PostgreSQL and MySQL. You need to configure the following environment variables in your .env file:

Locate the .env file inside the project folder “spurtcms-admin-app” and configure it with the details of newly created database to use either PostgreSQL or MySQL such as database type,name, user name, password etc





### Database Configuration

```
DATABASE_TYPE=postgres          # Use 'postgres' for PostgreSQL or 'mysql' for MySQL
DB_HOST=localhost               # Database host
DB_PORT=5432                    # Port for PostgreSQL (use 3306 for MySQL)
DB_DATABASE=your_database_name  # Name of your database
DB_USERNAME=your_username       # Database username
DB_PASSWORD=your_password       # Database password
DB_SSL_MODE=disable

```






Successful completion of this step completes the database configuration for spurtCMS Admin application.

### Storage Configuration

SpurtCMS supports multiple storage backends for file uploads:

1. **Local Storage**: Files are stored in the local filesystem
2. **AWS S3**: Files are stored in Amazon S3 buckets
3. **Azure Storage**: (Coming soon)

Configure your storage preferences in the `.env` file:

```
# Storage Configuration
# Options: local, aws, azure
STORAGE_TYPE = 'local'
STORAGE_LOCAL_PATH = './storage'

# AWS S3 Configuration
AWS_ACCESS_KEY_ID = ''
AWS_SECRET_ACCESS_KEY = ''
AWS_DEFAULT_REGION = ''
AWS_BUCKET = ''

# Azure Storage Configuration
AZURE_STORAGE_ACCOUNT = ''
AZURE_ACCOUNT_KEY = ''
AZURE_CONTAINER_NAME = ''
```

Successful completion of this step completes the storage configuration for spurtCMS Admin application.


### Step 4: Running the Project

Open the terminal within the project / cloned folder “spurtcms-admin”, and execute the following command:

```
go run main.go

```
This command initiates the spurtCMS Admin application, allowing you to begin your journey with this powerful content management system.

 

By following the steps outlined in this article, you have successfully set up spurtCMS Admin on your system. Ensure that all prerequisites are met and the configuration steps are accurately executed to enjoy a seamless experience with spurtCMS Admin application. Now you can explore the features and functionalities of spurtCMS Admin for efficient content management.


live demo of our intuitive Admin Panel .

```
Username : spurtcmsAdmin
Password : Admin@123
```


## 🤔 Support , Document and Help

spurtcms 4.8.2 is published to npm under the `@spurtcms/*` namespace.

You can find our extended documentation on our [https://spurtcms.com/documentation](https://spurtcms.com/documentation), but some quick links that might be helpful:

- Read [Technology](https://www.spurtcms.com/opensource-ecommerce-multivendor-nodejs-react-angular) to learn about our vision and what's in the box.

- Our [Discard](https://discord.com/invite/9TNgqUY24N) Questions, Live Discussions [spurtcms Support](https://picco.support).

- Some [Video](https://www.youtube.com/@spurtcms/videos) Video Tutorials 
- Every [Release](https://github.com/spurtcms/spurtcms-admin/releases) is documented on the Github Releases page.

🐞 If you spot a bug, please [submit a detailed issue](https://github.com/spurtcms/spurtcms-admin/issues/new), and wait for assistance.




## ❯ Maintainers
spurtcms is developed and maintain by [Piccosoft Software Labs India (P) Limited,](https://www.piccosoft.com).


## ❯ License

spurtcms is released under the [BSD-3-Clause License.](https://github.com/spurtcms/spurtcms/blob/master/LICENSE).




# SpurtCMS Docker 部署指南

本指南将帮助你使用 Docker 和 Docker Compose 部署 SpurtCMS，并配置为使用本地文件存储而非 AWS S3。

## 前提条件

- 安装 [Docker](https://docs.docker.com/get-docker/)
- 安装 [Docker Compose](https://docs.docker.com/compose/install/)

## 部署步骤

### 1. 准备目录结构

确保在项目根目录下创建以下目录：

```bash
mkdir -p db static_files
```

- `db` 目录将用于存储 PostgreSQL 数据库文件
- `static_files` 目录将用于存储上传的媒体文件和其他静态资源

### 2. 配置环境变量

项目中已包含 `.env.docker` 文件，其中包含了 Docker 环境下的默认配置。如需自定义，可以编辑此文件。

### 3. 构建并启动容器

在项目根目录下运行以下命令：

```bash
docker-compose up -d
```

这将构建 SpurtCMS 镜像并启动所有服务。首次启动可能需要几分钟时间。

### 4. 访问 SpurtCMS

启动完成后，可以通过以下地址访问 SpurtCMS：

- 管理面板: http://localhost:8082
- 预览: http://localhost:8083
- GraphQL API: http://localhost:8084

默认登录凭据：
- 用户名: spurtcmsAdmin
- 密码: Admin@123

## 文件存储配置

本配置已设置 SpurtCMS 使用本地文件存储而非 AWS S3。所有上传的文件将存储在 `static_files` 目录中，该目录已映射到容器内的 `/app/storage` 路径。

数据库初始化脚本 `init-local-storage.sql` 会自动配置 SpurtCMS 使用本地存储。

## 目录结构

- `Dockerfile`: 定义 SpurtCMS 应用的 Docker 镜像
- `docker-compose.yml`: 定义服务、网络和卷
- `.env.docker`: Docker 环境的环境变量配置
- `init-local-storage.sql`: 数据库初始化脚本，配置本地存储
- `db/`: PostgreSQL 数据库文件目录
- `static_files/`: 静态文件和上传媒体的存储目录

## 常见问题

### 如何查看日志？

```bash
docker-compose logs -f spurtcms
```

### 如何重启服务？

```bash
docker-compose restart spurtcms
```

### 如何完全重建环境？

```bash
docker-compose down
docker-compose up -d --build
```

注意：这不会删除数据库和静态文件，因为它们存储在宿主机上的映射目录中。

### 如何备份数据？

备份数据库和静态文件目录：

```bash
# 备份数据库
docker exec spurtcms_postgres pg_dump -U spurtcms spurtcms > spurtcms_backup.sql

# 备份静态文件
tar -czf static_files_backup.tar.gz static_files/
```

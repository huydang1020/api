# HuyShop API Gateway

Đây là API gateway/service chính cho dự án HuyShop. Nó hoạt động như một điểm truy cập chính cho các request từ client, định tuyến chúng đến các microservices backend phù hợp thông qua gRPC và xử lý các tác vụ chung như xác thực (authentication), caching, và upload file.

## Công Nghệ Sử Dụng (Tech Stack)
- **Ngôn ngữ**: Go 1.23+
- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Giao tiếp**: REST (Dành cho Frontend) & gRPC (Giao tiếp nội bộ giữa các Backend)
- **Caching**: Redis
- **Lưu trữ Media**: Cloudinary

## Yêu Cầu Hệ Thống (Prerequisites)
- Go 1.23 hoặc mới hơn
- Redis server
- Tài khoản Cloudinary (nếu cần xử lý upload hình ảnh)
- Các backend gRPC services đang chạy (User, Product, Permission, Voucher)

## Cấu Hình (Configuration)

Ứng dụng sử dụng các biến môi trường để cấu hình. Bạn có thể tạo file `.env` ở thư mục gốc để phát triển ở môi trường local (copy từ `.env.example` nếu có).

### Các Biến Môi Trường

| Biến | Mô tả | Mặc định |
|----------|-------------|---------|
| `PORT` | Cổng (port) mà HTTP server sẽ lắng nghe | `8080` |
| `JWT_SECRET_KEY` | Secret key dùng để ký và xác thực JWT tokens | *trống* |
| `ADMIN_ROLE` | Tên role được sử dụng cho quyền admin | `admin` |
| `PERM_GRPC_SERVER` | Địa chỉ của Permission gRPC service | `localhost:7000` |
| `USER_GRPC_SERVER` | Địa chỉ của User gRPC service | `localhost:6000` |
| `VOUCHER_GRPC_SERVER`| Địa chỉ của Voucher gRPC service | `localhost:4000` |
| `PRODUCT_GRPC_SERVER`| Địa chỉ của Product gRPC service | `localhost:8000` |
| `REDIS_ADDR` | Địa chỉ Redis server | `localhost:6379` |
| `REDIS_PASSWORD` | Mật khẩu Redis | *trống* |
| `REDIS_DB` | Số database của Redis | `0` |
| `CLOUDINARY_NAME` | Cloud name của Cloudinary | *trống* |
| `CLOUDINARY_API_KEY` | API key của Cloudinary | *trống* |
| `CLOUDINARY_SECRET` | API secret của Cloudinary | *trống* |

## Chạy Ứng Dụng (Local)

Để build và chạy ứng dụng ở môi trường local, bạn có thể sử dụng `Makefile` đã được cung cấp sẵn:

```bash
# Lệnh này sẽ build Go binary và chạy ứng dụng
make start
```

Hoặc, sử dụng các lệnh Go tiêu chuẩn:
```bash
go build -o api .
./api start
```

## Triển Khai Với Docker (Docker Deployment)

Ứng dụng có bao gồm một `Dockerfile` để dễ dàng đóng gói (containerization). `Dockerfile` này sử dụng image Alpine và yêu cầu file binary `api` phải được build trước khi chạy `docker build`.

1. **Build binary cho Linux:**
   `Makefile` chứa lệnh để biên dịch tĩnh Go binary cho môi trường Linux:
   ```bash
   make build
   ```

2. **Build Docker image:**
   ```bash
   docker build -t huyshop-api .
   ```

3. **Chạy container:**
   ```bash
   docker run -p 8080:8080 --env-file .env huyshop-api
   ```

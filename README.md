# miniproject-golang

project-name/
├── cmd/ # จุดเริ่มต้นของแอปพลิเคชัน (main.go)
│ └── main.go
├── internal/ # โค้ดภายในที่ไม่ควรใช้งานจากภายนอก
│ ├── adapters/ # Adapter สำหรับเชื่อมต่อโลกภายนอก
│ │ ├── http/ # ส่วนขาเข้า (Primary Adapter)
│ │ │ ├── handlers/ # เชื่อมต่อกับ framework เช่น Fiber, Gin
│ │ │ ├── routes/ # กำหนดเส้นทาง API
│ │ ├── repositories/ # ขาออก (Secondary Adapter) เช่น database
│ ├── core/ # Business Logic หลักของระบบ
│ │ ├── domain/ # โครงสร้างข้อมูลหลัก เช่น User, Product
│ │ ├── ports/ # Interface สำหรับเชื่อมต่อ Adapter
│ │ ├── service/ # การประมวลผลหลักตาม use case
├── pkg/ # โค้ดที่สามารถนำกลับมาใช้ภายนอกได้
│ ├── configs/ # ไฟล์ config เช่น config.yaml
│ ├── middleware/ # Middleware ต่าง ๆ เช่น auth, logging
│ ├── utils/ # ฟังก์ชันช่วยเหลือทั่วไป
├── docs/ # เอกสารประกอบ เช่น OpenAPI spec
├── go.mod # Go module file
├── go.sum
└── README.md

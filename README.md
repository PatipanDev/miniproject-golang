# miniproject-golang


# project map
project-name/
├── cmd/ # จุดเริ่มต้นของแอปพลิเคชัน (main.go)
│ └── main.go
├── internal/ # โค้ดภายในที่ไม่ควรใช้งานจากภายนอก
| ├── adapters/ # domain ย่อย เช่น user
| | ├── http/ # http ส่วน ขาเชื่อมต่อ primary 
| | | ├── handlers/ # ส่วน ของการเชื่อมต่อภายนอก เช่น fiber gin
| | | ├── routes/ # กำหนดเส้นทางที่เข้ามาเชื่อมต่อ
| | ├── repositories/ # ส่วนของ second ในการเชื่อมต่อฐานข้อมูล
| ├── core/ # core ในส่วนของ business logic
| | ├── domain/ # domain ย่อย เช่น user
| | ├── ports/ # ports กำหนด interface หรือ ช่องทางในการเชื่อมต่อ
| | ├── service/ # service ส่วนหลักของ business logic
├── pkg/ # โค้ดที่อาจถูกเรียกใช้ซ้ำภายนอก
| ├── configs/ # config files เช่น config.yaml
| ├── middleware/ # ตรวจสอบค่าต่างๆที่เข้ามา
| ├── utils/ # ส่วนของฟังชันก์ ต่างๆ
├── docs/ # เอกสารต่าง ๆ เช่น OpenAPI spec
├── go.mod # Go module file
├── go.sum
└── README.md



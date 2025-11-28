namespace go product

// 商品查询请求参数，对应 ProductRequestDTO
struct ProductRequest {
    1: optional i64 id,             // 商品ID
    2: optional string productName, // 商品名称
    3: optional string category,    // 商品分类（直接存储名称）
    4: optional double price,       // 商品售价
    5: optional string description, // 商品简介
    6: optional i8 status,          // 状态：0-下架，1-上架
    7: optional string createTime   // 创建时间（如 "2000-01-01 12:00:00"）
}

// 商品响应信息，对应 ProductResponseDTO
struct ProductResponse {
    1: optional i64 id,             // 商品ID
    2: optional string productName, // 商品名称
    3: optional string category,    // 商品分类（直接存储名称）
    4: optional double price,       // 商品售价
    5: optional string description, // 商品简介
    6: optional i8 status,          // 状态：0-下架，1-上架
    7: optional string createTime,  // 创建时间（如 "2000-01-01 12:00:00"）
    8: optional i32 stockNum        // 库存数量
}

// 业务异常
exception ProductException {
    1: required string message,     // 错误详情
    2: optional i32 code            // 错误码
}

// 商品服务接口
service ProductService {
    // 根据ID获取商品详情（含库存）
    ProductResponse getProductById(1: ProductRequest request) throws (1: ProductException ex),
}

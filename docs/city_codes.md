# 城市编码映射表

## 城市编码说明

本系统使用标准的行政区划代码作为城市编码，前端可以根据需要将编码转换为城市名称。

## 主要城市编码

| 城市编码 | 城市名称 | 省份 | 区域 |
|---------|---------|------|------|
| 110000 | 北京市 | 北京市 | 华北 |
| 310000 | 上海市 | 上海市 | 华东 |
| 440100 | 广州市 | 广东省 | 华南 |
| 440300 | 深圳市 | 广东省 | 华南 |
| 330100 | 杭州市 | 浙江省 | 华东 |
| 320100 | 南京市 | 江苏省 | 华东 |
| 510100 | 成都市 | 四川省 | 西南 |
| 420100 | 武汉市 | 湖北省 | 华中 |
| 610100 | 西安市 | 陕西省 | 西北 |
| 500000 | 重庆市 | 重庆市 | 西南 |

## 前端使用示例

### JavaScript 城市映射
```javascript
const cityMapping = {
  '110000': { name: '北京市', province: '北京市', region: '华北' },
  '310000': { name: '上海市', province: '上海市', region: '华东' },
  '440100': { name: '广州市', province: '广东省', region: '华南' },
  '440300': { name: '深圳市', province: '广东省', region: '华南' },
  '330100': { name: '杭州市', province: '浙江省', region: '华东' },
  '320100': { name: '南京市', province: '江苏省', region: '华东' },
  '510100': { name: '成都市', province: '四川省', region: '西南' },
  '420100': { name: '武汉市', province: '湖北省', region: '华中' },
  '610100': { name: '西安市', province: '陕西省', region: '西北' },
  '500000': { name: '重庆市', province: '重庆市', region: '西南' }
};

// 获取城市名称
function getCityName(cityCode) {
  return cityMapping[cityCode]?.name || cityCode;
}

// 获取城市信息
function getCityInfo(cityCode) {
  return cityMapping[cityCode] || { name: cityCode, province: '', region: '' };
}
```

### React 组件示例
```jsx
import React from 'react';

const CityDisplay = ({ cityCode }) => {
  const cityInfo = getCityInfo(cityCode);
  
  return (
    <span>{cityInfo.name}</span>
  );
};

// 城市选择组件
const CitySelect = ({ value, onChange }) => {
  return (
    <select value={value} onChange={(e) => onChange(e.target.value)}>
      {Object.entries(cityMapping).map(([code, info]) => (
        <option key={code} value={code}>
          {info.name}
        </option>
      ))}
    </select>
  );
};
```

## 数据格式

### 客户授权记录表 (customer_authorization_record)
- 字段名: `city_code`
- 类型: `varchar(10)`
- 说明: 存储城市编码，如 '110000' 表示北京

### 合作门店表 (cooperation_store)
- 字段名: `cooperation_city_code`
- 类型: `varchar(10)`
- 说明: 存储合作城市编码，如 '310000' 表示上海

## 注意事项

1. **编码标准**: 使用国家标准的行政区划代码
2. **前端处理**: 前端负责将编码转换为用户友好的城市名称
3. **数据一致性**: 后端只存储编码，确保数据一致性
4. **扩展性**: 可以根据需要添加更多城市编码

## 常用城市编码查询

可以通过以下方式获取完整的城市编码列表：

1. **国家统计局**: 提供最新的行政区划代码
2. **第三方库**: 如 `china-area-data` 等 npm 包
3. **API 服务**: 使用城市数据 API 服务

## 更新记录

- 2024-01-13: 初始版本，支持主要一二线城市
- 后续可根据业务需要扩展更多城市

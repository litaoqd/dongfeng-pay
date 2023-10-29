<?php
// 允许跨域请求
header('Access-Control-Allow-Origin: *');
header('Content-Type: application/json');

// API URL 和认证信息
$apiUrl = 'https://api.onepayph.com/api/v1/payment_intent/ref_trade_id/';
$bearerToken = 'Bearer sk_test_2WoGEvALkxMmU0oy6WZsHzQXTpO';

// 获取前端传递的商户订单号
$merchantOrderId = $_GET['merchantOrderId'] ?? '';

if (!$merchantOrderId) {
    echo json_encode(['error' => 'No merchant order ID provided']);
    exit;
}

// 构建完整的 API 请求 URL
$url = $apiUrl . urlencode($merchantOrderId);

echo $url;

// 使用 cURL 发送请求
$ch = curl_init($url);
curl_setopt($ch, CURLOPT_HTTPHEADER, ['Authorization: ' . $bearerToken]);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);

$response = curl_exec($ch);
curl_close($ch);

// 输出响应
echo $response;

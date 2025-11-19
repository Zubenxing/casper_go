# 测试 Azure 账单对比 Webhook
$webhookUrl = "http://localhost:5678/webhook/azure-bill"

Write-Host "测试 Webhook: $webhookUrl" -ForegroundColor Cyan

# 创建测试文件（空的也行）
$testFile1 = [System.IO.Path]::GetTempFileName()
$testFile2 = [System.IO.Path]::GetTempFileName()

Write-Host "创建测试文件..." -ForegroundColor Yellow

# 发送 POST 请求
try {
    $response = Invoke-WebRequest -Uri $webhookUrl -Method POST -Form @{
        current_month = Get-Item $testFile1
        last_month = Get-Item $testFile2
    } -UseBasicParsing
    
    Write-Host "`n✅ 请求成功!" -ForegroundColor Green
    Write-Host "Status Code: $($response.StatusCode)" -ForegroundColor Green
    Write-Host "Content-Type: $($response.Headers['Content-Type'])" -ForegroundColor Cyan
    Write-Host "Content-Length: $($response.RawContentLength)" -ForegroundColor Cyan
    
    # 检查是否是文件
    $contentType = $response.Headers['Content-Type']
    if ($contentType -like "*excel*" -or $contentType -like "*sheet*") {
        Write-Host "`n✅ 返回的是 Excel 文件！" -ForegroundColor Green
        
        # 保存文件
        $outputFile = "test_output.xlsx"
        [System.IO.File]::WriteAllBytes($outputFile, $response.Content)
        Write-Host "文件已保存到: $outputFile" -ForegroundColor Green
    } else {
        Write-Host "`n⚠️ 返回的不是 Excel 文件" -ForegroundColor Yellow
        Write-Host "响应内容: $($response.Content)" -ForegroundColor Gray
    }
} catch {
    Write-Host "`n❌ 请求失败!" -ForegroundColor Red
    Write-Host "错误: $_" -ForegroundColor Red
} finally {
    # 清理测试文件
    Remove-Item $testFile1 -ErrorAction SilentlyContinue
    Remove-Item $testFile2 -ErrorAction SilentlyContinue
}

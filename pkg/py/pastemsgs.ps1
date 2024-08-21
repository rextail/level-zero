# Определяем путь к файлу
$filePath = "generated_data.json"

# Загружаем JSON-файл и разбиваем его на отдельные сообщения
$messages = Get-Content -Raw -Path $filePath | ConvertFrom-Json

# Цикл по каждому сообщению
foreach ($message in $messages) {
    # Преобразуем сообщение обратно в JSON
    $jsonMessage = $message | ConvertTo-Json -Depth 10

    # Отправляем сообщение в NATS через стандартный ввод
    $jsonMessage | nats pub orders.new

    # Пауза между отправками (50 миллисекунд)
    Start-Sleep -Milliseconds 50
}

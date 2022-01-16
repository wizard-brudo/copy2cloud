echo Компиляция прогграмы
go build -o copy2cloud main.go client.go
echo Создание директории в opt
sudo mkdir /opt/copy2cloud
echo Копирование бинарника в /opt/copy2cloud
sudo mv copy2cloud /opt/copy2cloud
echo Копирование шаблонов в /opt/copy2cloud
sudo cp -r oauth2/templates /opt/copy2cloud
echo Добавление прогграмы в .bashrc
echo "export PATH=\$PATH:/opt/copy2cloud" >> ~/.bashrc

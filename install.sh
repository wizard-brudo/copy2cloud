echo Компиляция прогграмы
go build -o copy2cloud main.go client.go
echo Создание директории в opt
mkdir ~/copy2cloud
echo Копирование бинарника в ~/copy2cloud
mv copy2cloud ~/copy2cloud
echo Копирование шаблонов в ~/copy2cloud
cp -r oauth2/templates ~/copy2cloud
echo Добавление прогграмы в .bashrc
echo "export PATH=\$PATH:~/copy2cloud" >> ~/.bashrc

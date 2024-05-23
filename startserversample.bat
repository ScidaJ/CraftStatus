:: This will automatically restart your server when sent the /stop command. 

:START
:: Place code which starts your server here
echo Restarting automatically in 10 seconds (press Ctrl + C to cancel)
timeout /t 10 /nobreak > NUL
goto:START
@REM install this service in windows
@REM
@REM
SET SERVICENAME=rift
SET APP=rift
SET NSSM=C:\Nssm\win64\nssm.exe
FOR /F "delims=" %%i IN ('dir %APP%.cmd /b/s') DO set STARTUPSCRIPT=%%i

%NSSM% install %SERVICENAME% %STARTUPSCRIPT%
%NSSM% set %SERVICENAME% AppEnvironmentExtra JAVA_HOME=%JAVA_HOME%
%NSSM% set %SERVICENAME% AppDirectory %~dp0
%NSSM% start %SERVICENAME%

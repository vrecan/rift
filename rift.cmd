@REM       app="metrics"
@REM       java=/usr/java/default/bin/java
@REM       jarPath=`find $PWD -name $app-*.jar -print -quit`
@REM       logback="${app}-logback-persistent.xml"
@REM       if [ ! -e $logback ]
@REM       then
@REM          unzip -q -c $jarPath logback-persistent.xml > $logback
@REM          fi
@REM          exec /bin/sh -c "$java -Dlogback.configurationFile=$logback -jar $jarPath"
SET APP=rift
./rift.exe
goto :finally 

:err
echo JAVA_HOME environment variable must be set!
pause

REM -----------------------------------------------------------------------------
:finally

ENDLOCAL

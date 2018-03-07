# For windows developers using powershell
echo "DROP DATABASE caraway" | psql -U postgres
cat $Env:GOPATH\src\395_Project_2018\database\dbinit.sql | psql -U postgres
cat $Env:GOPATH\src\395_Project_2018\database\data.sql | psql -U caraway
cat $Env:GOPATH\src\395_Project_2018\database\bookingdata.sql | psql -U caraway
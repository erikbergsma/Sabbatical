var app = angular.module('myApp', []);
app.controller('customersCtrl', function($scope, $http) {
  //get data from API
  $http.get("/api/list")
  .then(function (response) {
    $scope.names = response.data;

    angular.forEach($scope.names, function (item) {
      if (item.Enabled){
        item.selectedOption = {id: 0, name: true}
      } else{
        item.selectedOption = {id: 1, name: false}
      }
    })
  });

  $scope.create = function(array){
    var name = $('#newName').val()
    var enabledval = $('#newEnabledSelect').val()

    if (enabledval == "true"){
      var selectedOption = {id: 0, name: true}
      var enabled = true
    } else {
      var selectedOption = {id: 1, name: false}
      var enabled = false
    }

    var newItem = {"Name": name, "Enabled": enabled, "selectedOption": selectedOption}

    $http.post('/api/create', JSON.stringify(newItem)).then(function (response) {
      if (response.data)
        //same as allways
        $scope.msg = "Post Data Submitted Successfully!"
        $scope.timestamp = getTime();

        //the id is generated at the backend
        newItem["ID"] = response.data["ID"]

        //give angularjs a new entry in the data, so it can update the HTML
        array.push(newItem);

        //clean the input html element
        $('#newName').val("")

      }, function (response) {
        userFeedback();
    });
  }

  //onclick handler
  $scope.remove = function(array, index){
    //Call the services
    $http.post('/api/delete', JSON.stringify(array[index])).then(function (response) {
      if (response.data)
        $scope.msg = "Post Data Submitted Successfully!"
        $scope.timestamp = getTime();
        array.splice(index, 1);
      }, function (response) {
        userFeedback();
    });
  }

  $scope.update = function(array, index, id){
    var name = $('#' + id + '_name').val()
    var enabledval = $('#' + id + '_enabled').val()
    var enabled = $scope.data.availableOptions[enabledval].name

    // i can just pass the index and the array, so i dont need to do a for loop ^^
    array[index].Name = name;
    array[index].Enabled = enabled;

    //Call the services
    $http.post('/api/update', JSON.stringify(array[index])).then(function (response) {
      if (response.data)
        $scope.msg = "Post Data Submitted Successfully!"
        $scope.timestamp = getTime();
      }, function (response) {
        userFeedback();
    });
  }

  function userFeedback(response){
    $scope.msg = response.xhrStatus;
    $scope.timestamp = getTime();
    $scope.statusval = response.status;
    $scope.statustext = response.statusText;
    $scope.headers = response.headers();
  }

  $scope.data = {
    availableOptions: [
      {id: 0, name: true},
      {id: 1, name: false}
    ]
  };
});

function getTime(){
  var today = new Date();
  var date = today.getDate()+'-'+(today.getMonth()+1)+'-'+today.getFullYear();
  var time = today.getHours() + ":" + today.getMinutes() + ":" + today.getSeconds();
  return date+' '+time;
}

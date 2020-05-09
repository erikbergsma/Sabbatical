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

'use strict';

var app = angular.module('application', []);

app.controller('AppCtrl', function($scope, appFactory){
    $("#success_init").hide();
    $("#success_invoke").hide();
    $("#success_query").hide();
    $("#success_qurey_all").hide();
    $("#success_delete").hide();

    $("#success_createLoan").hide();
    $("#success_approveLoan").hide();
    $("#success_deleteLoan").hide();
    $("#success_updateLoan").hide();
    $("#success_queryLoan").hide();
    $("#success_queryAllLoans").hide();

    // 기존 ABStore 기능
    $scope.initAB = function(){
        appFactory.initAB($scope.abstore, function(data){
            if(data == "Success")
                $scope.init_ab = "success";
            $("#success_init").show();
        });
    }

    $scope.invokeAB = function(){
        appFactory.invokeAB($scope.abstore, function(data){
            if(data == "Success")
                $scope.invoke_ab = "success";
            $("#success_invoke").show();
        });
    }

    $scope.queryAB = function(){
        appFactory.queryAB($scope.walletid, function(data){
            $scope.query_ab = data;
            $("#success_qurey").show();
        });
    }

    $scope.queryAll = function(){
        appFactory.queryAll(function(data){
            $scope.query_all = data;
            $("#success_qurey_all").show();
        });
    }

    $scope.deleteAB = function(){
        appFactory.deleteAB($scope.abstore, function(data){
            if(data == "Success")
                $scope.delete_ab = "success";
            $("#success_delete").show();
        });
    }

    // 깐부 대출 Loan 기능
    $scope.createLoan = function(){
        appFactory.createLoan($scope.loaninfo, function(data){
            if(data == "Success")
                $scope.create_loan = "Loan request created!";
            $("#success_createLoan").show();
        });
    }

    $scope.approveLoan = function(){
        appFactory.approveLoan($scope.loaninfo, function(data){
            if(data == "Success")
                $scope.approve_loan = "Loan request approved!";
            $("#success_approveLoan").show();
        });
    }

    $scope.deleteLoan = function(){
        appFactory.deleteLoan($scope.loaninfo, function(data){
            if(data == "Success")
                $scope.delete_loan = "Loan request deleted!";
            $("#success_deleteLoan").show();
        });
    }

    $scope.updateLoan = function(){
        appFactory.updateLoan($scope.loaninfo, function(data){
            if(data == "Success")
                $scope.update_loan = "Loan request updated!";
            $("#success_updateLoan").show();
        });
    }

    $scope.queryLoan = function(){
        appFactory.queryLoan($scope.loaninfo.loanid, function(data){
            $scope.query_loan = data;
            $("#success_queryLoan").show();
        });
    }

    $scope.queryAllLoans = function(){
        appFactory.queryAllLoans(function(data){
            $scope.query_all_loans = data;
            $("#success_queryAllLoans").show();
        });
    }
});

app.factory('appFactory', function($http){
    var factory = {};

    // 기존 ABStore API
    factory.initAB = function(data, callback){
        $http.get('/init?a='+data.a+'&aval='+data.aval+'&b='+data.b+'&bval='+data.bval)
            .success(function(output){
                callback(output)
            });
    }

    factory.invokeAB = function(data, callback){
        $http.get('/invoke?a='+data.a+'&b='+data.b+'&value='+data.value)
            .success(function(output){
                callback(output)
            });
    }

    factory.queryAB = function(a, callback){
        $http.get('/query?name='+a)
            .success(function(output){
                callback(output)
            });
    }

    factory.queryAll = function(callback){
        $http.get('/queryAll')
            .success(function(output){
                callback(output)
            });
    }

    factory.deleteAB = function(data, callback){
        $http.get('/delete?name='+data.a)
            .success(function(output){
                callback(output)
            });
    }

    // 깐부 대출 Loan API
    factory.createLoan = function(data, callback){
        $http.get('/createLoan?id='+data.loanid+'&requester='+data.requester+'&amount='+data.amount+'&durationDays='+data.durationDays)
            .success(function(output){
                callback(output)
            });
    }

    factory.approveLoan = function(data, callback){
        $http.get('/approveLoan?id='+data.loanid+'&provider='+data.provider)
            .success(function(output){
                callback(output)
            });
    }

    factory.deleteLoan = function(data, callback){
        $http.get('/deleteLoan?id='+data.loanid)
            .success(function(output){
                callback(output)
            });
    }

    factory.updateLoan = function(data, callback){
        $http.get('/updateLoan?id='+data.loanid+'&newAmount='+data.amount+'&newDurationDays='+data.durationDays)
            .success(function(output){
                callback(output)
            });
    }

    factory.queryLoan = function(id, callback){
        $http.get('/queryLoan?id='+id)
            .success(function(output){
                callback(output)
            });
    }

    factory.queryAllLoans = function(callback){
        $http.get('/queryAllLoans')
            .success(function(output){
                callback(output)
            });
    }

    return factory;
});

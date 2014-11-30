var captainWebhookApp = angular.module('captainWebhookApp', ['ui.bootstrap', 'ngRoute', 'ngDragDrop', 'frapontillo.bootstrap-switch'])
	.config(['$routeProvider', '$locationProvider', 
		function($routeProvider, $locationProvider){
			$routeProvider
				.when('/transformer/:id', {
		          templateUrl: '/html/transformer_config.html',
		          controller: 'captainWebhookTransformerEdit',
		        })
		        .when('/', {
		          templateUrl: '/html/transformer_home.html',
		          controller: 'captainWebhookTransformerHome',
		        });
		    $locationProvider.html5Mode(true);
		}]);

captainWebhookApp.filter('num', function() {
    return function(input) {
      return parseInt(input);
    }
});

captainWebhookApp.controller('captainWebhook', ['$scope', 'Transformers', '$route', '$routeParams', '$location', function($scope, Transformers, $route, $routeParams, $location){
	$scope.transformers = Transformers;
	this.$route = $route;
    this.$location = $location;
    this.$routeParams = $routeParams;
	Transformers.get();

	$scope.activeTransformerName = function(){
		if(angular.isDefined($scope.transformers.activeTransformer)){
			return $scope.transformers.activeTransformer.Name;
		}
		return 'Configuration';
	}
}]);

captainWebhookApp.directive('sortableObjectTransformations', [function(){
	return {
		link : function(scope, element, attrs){
			$(element).sortable(
				{
					placeholder: "sortable-btn placeholder-btn btn",
					forcePlaceholderSize : true
				}
			);
			$(element).on( "sortstop", function( ev, ui ) {
				var arr = $(element).sortable( "toArray" );
				scope.updateObjectTransformationWeights(arr);				
			});
		}
	}
}]);

captainWebhookApp.controller('captainWebhookTransformerHome', ['$scope', 'Transformers', '$location', function($scope, Transformers, $location){
	Transformers.get();
	$scope.transformers = Transformers;
	$scope.newTransformerName = '';

	$scope.addTransformer = function(){
		var newTrans = {Name : $scope.newTransformerName};
		Transformers.createTransformer(newTrans, function(resp){
			$location.url('/transformer/' + resp.Id);
		});
	}

	$scope.deleteTransformer = function(transId){
		$scope.transformers.deleteTransformer(transId);
		if(!$scope.$$phase){
			$scope.$apply();
		}
	}
}]);

captainWebhookApp.controller('captainWebhookTransformerEdit', ['$routeParams', '$scope', 'Transformers', '$location', function($routeParams, $scope, Transformers, $location){
	$scope.transformer;
	$scope.messages;	
	$scope.transformerId = $routeParams.id;
	$scope.activeMessage;
	$scope.activeRelId;
	$scope.sortableRelIdList;
	$scope.recalculating = false;
	$scope.transformationText;
	$scope.transformerAddress = $location.protocol() + '://' + $location.host() + ':' + $location.port() + '/receive/' + $scope.transformerId;

	Transformers.getOne($scope.transformerId, function(transformer){
		$scope.transformer = transformer;
		Transformers.activeTransformer = transformer;

		Transformers.getMessages($scope.transformerId, function(messages){
			$scope.messages = messages;
			$scope.activeMessage = messages[0];
			if(angular.isUndefined($scope.activeMessage)){
				return;
			}
			$scope.recalculateTransformation($scope.activeMessage.Id);
		});
	});

	$scope.saveTransformer = function(){
		Transformers.saveTransformer($scope.transformer);
		$scope.recalculateTransformation($scope.activeMessage.Id);
	}
	// Takes an array of Id's and updates the weight of the ObjectTransformation objects
	$scope.updateObjectTransformationWeights = function(idArr){
		for(var i = 0; i < idArr.length; i++){
			var j = $scope.indexOfOjbectTransformation(idArr[i]);
			$scope.transformer.ObjectTransformation[j].Weight = i;			
		}
		$scope.recalculateTransformation($scope.activeMessage.Id);
	}

	$scope.indexOfOjbectTransformation = function(relId){
		for(i =0; i < $scope.transformer.ObjectTransformation.length; i++){
			if($scope.transformer.ObjectTransformation[i].Rel_id == relId){
				return i;
			}
		}
		return -1;
	}

	$scope.setActiveMessageObject = function(relId){
		$scope.activeRelId = relId;
		$scope.recalculateTransformation($scope.activeMessage.Id);
		var indexOfOT = $scope.indexOfOjbectTransformation(relId);
		if(indexOfOT == -1){
			$scope.transformer.ObjectTransformation.push({
				Template : '',
				Rel_id : relId,
				Weight : $scope.transformer.ObjectTransformation.length
			});
		}
	}

	$scope.recalculateTransformation = function(messageId){
		if($scope.recalculating){
			return;
		}
		$scope.recalculating = true;
		Transformers.translateMessage(messageId, $scope.transformer, function(translation){
			$scope.recalculating = false;
			$scope.transformationText = translation;
		});
	}

	$scope.addProperty = function(relId, key){
		var indexOfOT = $scope.indexOfOjbectTransformation(relId);
		$scope.transformer.ObjectTransformation[indexOfOT].Template = $scope.transformer.ObjectTransformation[indexOfOT].Template + " {{." + key + "}}"
	}

	$scope.setActiveMessage = function(message){		
		$scope.activeMessage = message;
		$scope.recalculateTransformation($scope.activeMessage.Id);
	}
}]);
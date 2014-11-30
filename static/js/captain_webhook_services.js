captainWebhookApp.service('Transformers', ["$http", "$rootScope", function($http, $rootScope){
	var Transformers = {};
	Transformers.transformers;
	Transformers.activeTransformer;
	Transformers.get = function(){
		if(typeof Transformers.transformers == 'undefined'){
			$http.get('/api/transformers').success(function(response){
				Transformers.transformers = response;
			});
		}
	}

	Transformers.getOne = function(id, callback){
		$http.get('/api/transformers/' + id).success(function(response){
			callback(response);
		});
	}

	Transformers.getMessages = function(id, callback){
		$http.get('/api/transformers/' + id + '/messages').success(function(response){
			callback(response);
		});
	}

	Transformers.saveTransformer = function(transformer, callback){
		$http.put('/api/transformers/' + transformer.Id, transformer).success(function(response){
			if(angular.isDefined(callback)){
				callback(response);
			}
		});
	}

	Transformers.translateMessage = function(messageId, transformer, callback){
		$http.post('/api/transform_message/' + messageId, transformer, {
			// Angular getting too smart for its own good
			transformResponse : function(data, headers){
				return data;
			}
		}).success(function(response){
			callback(response);
		})
	}

	return Transformers;
}]);
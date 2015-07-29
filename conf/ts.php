{phptag}
if(!class_exists('TS')){
	class TS {
		static $token = '{token}';
		static $url = 'http://{domain}/ts/post';
		static $data = array();
		static function T(){
			self::_A(debug_backtrace());
		}

		static function S(){
			$post = array(
			    'type' => 'php',
				'data' => self::_C(self::$data),
				'token' => self::$token,
				'time' => self::_time()
			);
			$ch = curl_init();
			curl_setopt($ch, CURLOPT_URL, self::$url);
			curl_setopt($ch, CURLOPT_POST, true);
			curl_setopt($ch, CURLOPT_TIMEOUT, 30 * 1000);
			curl_setopt($ch, CURLOPT_POSTFIELDS, $post);
			curl_exec($ch);
			curl_close($ch);
		}

		static function T_S(){
			self::$data = array();
			self::_A(debug_backtrace());
			self::S();
		}

		static function _A($debug){
			$data = array(
				'debug' => $debug,
				'time' => self::_time()
			);
			self::$data[] = $data;
		}

		static function _C($array){
			return urlencode(json_encode($array));
		}

		static function _time(){
		    return floor(microtime(true) * 1000);
		}
	}
	return true;
}

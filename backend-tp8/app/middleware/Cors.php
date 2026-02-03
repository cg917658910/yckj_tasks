<?php

declare(strict_types=1);

namespace app\middleware;

use think\Request;
use think\Response;

class Cors
{
    public function handle(Request $request, \Closure $next): Response
    {
        if ($request->isOptions()) {
            return response('', 204)->header($this->headers());
        }
        
        return $next($request)->header($this->headers());
    }

    private function headers(): array
    {
        return [
            'Access-Control-Allow-Origin' => '*',
            'Access-Control-Allow-Methods' => 'GET,POST,PUT,PATCH,DELETE,OPTIONS',
            'Access-Control-Allow-Headers' => 'Authorization,Content-Type',
        ];
    }
}

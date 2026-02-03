<?php

declare(strict_types=1);

namespace app\controller\admin;

use app\controller\BaseController;
use app\service\AuthService;

class AuthController extends BaseController
{
    public function login()
    {
        $username = trim((string) $this->request->post('username', ''));
        $password = (string) $this->request->post('password', '');

        if ($username === '' || $password === '') {
            return json_error('参数不完整');
        }

        try {
            $data = AuthService::adminLogin($username, $password);
        } catch (\Throwable $e) {
            return json_error($e->getMessage());
        }

        return json_success($data);
    }

    public function logout()
    {
        return json_success();
    }
}

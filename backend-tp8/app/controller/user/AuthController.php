<?php

declare(strict_types=1);

namespace app\controller\user;

use app\controller\BaseController;
use app\service\AuthService;

class AuthController extends BaseController
{
    public function register()
    {
        $data = $this->request->only(['username', 'password'], 'post');
        try {
            validate(\app\validate\UserRegisterValidate::class)->check($data);
        } catch (\think\exception\ValidateException $e) {
            return json_error($e->getError());
        }

        $username = trim((string) $data['username']);
        $password = (string) $data['password'];

        try {
            $data = AuthService::userRegister($username, $password);
        } catch (\Throwable $e) {
            return json_error($e->getMessage());
        }

        return json_success($data);
    }

    public function login()
    {
        $data = $this->request->only(['username', 'password'], 'post');
        try {
            validate(\app\validate\UserLoginValidate::class)->check($data);
        } catch (\think\exception\ValidateException $e) {
            return json_error($e->getError());
        }

        $username = trim((string) $data['username']);
        $password = (string) $data['password'];

        try {
            $data = AuthService::userLogin($username, $password);
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

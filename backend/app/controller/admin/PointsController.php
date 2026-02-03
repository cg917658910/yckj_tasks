<?php

declare(strict_types=1);

namespace app\controller\admin;

use app\controller\BaseController;
use app\service\PointsService;

class PointsController extends BaseController
{
    public function rules()
    {
        $rule = PointsService::getRules();
        return json_success($rule->toArray());
    }

    public function updateRules()
    {
        $adminId = (int) ($this->request->admin['admin_id'] ?? 0);
        $data = [
            'exchange_rate' => (int) $this->request->post('exchange_rate', 10),
            'min_withdraw_amount' => (float) $this->request->post('min_withdraw_amount', 10),
            'register_bonus_points' => (int) $this->request->post('register_bonus_points', 10),
        ];

        try {
            $rule = PointsService::updateRules($adminId, $data);
        } catch (\Throwable $e) {
            return json_error($e->getMessage());
        }

        return json_success($rule->toArray());
    }
}

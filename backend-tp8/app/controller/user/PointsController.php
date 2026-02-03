<?php

declare(strict_types=1);

namespace app\controller\user;

use app\controller\BaseController;
use think\facade\Db;

class PointsController extends BaseController
{
    public function logs()
    {
        $userId = (int) ($this->request->user['user_id'] ?? 0);
        $page = (int) $this->request->get('page', 1);
        $pageSize = (int) $this->request->get('page_size', 10);

        $query = Db::name('points_logs')->alias('l')
            ->field('l.*')
            ->where('l.user_id', $userId)
            ->order('l.id', 'desc');

        $pager = $query->paginate(['list_rows' => $pageSize, 'page' => $page])->toArray();

        return json_success([
            'list' => $pager['data'] ?? [],
            'total' => (int) ($pager['total'] ?? 0),
            'page' => (int) ($pager['current_page'] ?? $page),
            'page_size' => (int) ($pager['per_page'] ?? $pageSize),
        ]);
    }
}

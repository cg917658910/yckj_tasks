<?php

declare(strict_types=1);

namespace app\controller\admin;

use app\controller\BaseController;
use app\model\Task;
use app\service\TaskService;

class TaskController extends BaseController
{
    public function index()
    {
        $status = $this->request->get('status');
        $query = Task::order('id', 'desc');
        if ($status !== null && $status !== '') {
            $query->where('status', (int) $status);
        }
        $list = $query->select()->toArray();
        return json_success(['list' => $list]);
    }

    public function create()
    {
        $data = [
            'title' => trim((string) $this->request->post('title', '')),
            'summary' => trim((string) $this->request->post('summary', '')),
            'detail' => (string) $this->request->post('detail', ''),
            'doc_url' => (string) $this->request->post('doc_url', ''),
            'reward_points' => (int) $this->request->post('reward_points', 0),
            'status' => TaskService::TASK_STATUS_ONLINE,
        ];

        if ($data['title'] === '' || $data['summary'] === '' || $data['detail'] === '') {
            return json_error('参数不完整');
        }

        $task = Task::create($data);
        return json_success($task->toArray());
    }

    public function update(int $id)
    {
        $task = Task::where('id', $id)->find();
        if (!$task) {
            return json_error('任务不存在');
        }

        $task->title = trim((string) $this->request->post('title', $task->title));
        $task->summary = trim((string) $this->request->post('summary', $task->summary));
        $task->detail = (string) $this->request->post('detail', $task->detail);
        $task->doc_url = (string) $this->request->post('doc_url', $task->doc_url);
        $task->reward_points = (int) $this->request->post('reward_points', $task->reward_points);
        $task->status = (int) $this->request->post('status', $task->status);
        $task->save();

        return json_success($task->toArray());
    }

    public function off(int $id)
    {
        $task = Task::where('id', $id)->find();
        if (!$task) {
            return json_error('任务不存在');
        }

        $task->status = TaskService::TASK_STATUS_OFFLINE;
        $task->save();

        return json_success();
    }
}

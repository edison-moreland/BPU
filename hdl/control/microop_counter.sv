module microop_counter(
    input logic clk, inst_done,
    output logic [0:2] count,
);
    logic should_rst = 0;
    always_ff @(posedge clk, posedge inst_done) begin
        unique if (inst_done)
            should_rst <= 1;
        else if (clk) begin
            if (should_rst) begin
                count <= 0;
                should_rst <= 0;
            end
            else
                count <= count + 1;
        end
    end

endmodule
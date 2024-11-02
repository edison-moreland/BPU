module tflop(input  logic clk,
             input  logic t,
             output logic q);

    always @(posedge clk) begin
        if(t == 0)
            q <= q;
        else
            q = ~q;
    end
endmodule